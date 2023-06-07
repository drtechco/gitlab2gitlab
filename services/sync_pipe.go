package services

import (
	"bytes"
	"context"
	"drtech.co/gl2gl/core/configs"
	"drtech.co/gl2gl/orm"
	"errors"
	"fmt"
	"github.com/ahmetb/go-linq/v3"
	"github.com/ldez/go-git-cmd-wrapper/v2/clone"
	"github.com/ldez/go-git-cmd-wrapper/v2/fetch"
	git "github.com/ldez/go-git-cmd-wrapper/v2/git"
	"github.com/ldez/go-git-cmd-wrapper/v2/push"
	"github.com/ldez/go-git-cmd-wrapper/v2/types"
	"github.com/sirupsen/logrus"
	"github.com/xanzy/go-gitlab"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

type SyncPipe struct {
	ConfigId        int32
	FromAddress     string
	FromAccessToken string
	FromClient      *gitlab.Client
	ToAddress       string
	ToAccessToken   string
	ToClient        *gitlab.Client
	Status          SyncPipeStatus

	logger *logrus.Entry
}

func (p *SyncPipe) Stop() error {
	return nil
}

func (p *SyncPipe) Run() error {
	p.logger = logrus.WithField("Name", "SyncPipe")
	p.UpdateStatus(SyncPipeStatusIniting)
	fromClient, err := gitlab.NewClient(p.FromAccessToken, gitlab.WithBaseURL(
		fmt.Sprintf("%s/api/v4", p.FromAddress),
	))
	if err != nil {
		p.UpdateStatus(SyncPipeStatusFromClientInitError)
		return err
	}
	toClient, err := gitlab.NewClient(p.ToAccessToken, gitlab.WithBaseURL(
		fmt.Sprintf("%s/api/v4", p.ToAddress),
	))
	if err != nil {
		p.UpdateStatus(SyncPipeStatusToClientInitError)
		return err
	}
	p.UpdateStatus(SyncPipeStatusInitOk)
	p.FromClient = fromClient
	p.ToClient = toClient
	p.Sync()
	return nil
}

func (p *SyncPipe) UpdateStatus(status SyncPipeStatus) {
	fromToConfigM := orm.DbQuery().FromToConfig
	_, err := fromToConfigM.WithContext(context.Background()).
		Where(fromToConfigM.ID.Eq(p.ConfigId)).Update(fromToConfigM.Status, status)
	if err != nil {
		p.logger.Error("UpdateStatus:", err)
	}
}

func (p *SyncPipe) Sync() {
	for {
		err := p.SyncUser()
		if err != nil {
			p.logger.Error("SyncUser:", err)
		}
		err = p.SyncGroup(nil, nil)
		if err != nil {
			p.logger.Error("SyncGroup:", err)
		}
		time.Sleep(time.Second * 120)
	}
}

func (p *SyncPipe) SyncGroup(parentFrom *gitlab.Group, parentTo *gitlab.Group) error {
	if parentFrom == nil {
		p.logger.Debug("开始同步Root Group")
		topLevelOnly := true
		p.UpdateStatus(SyncPipeStatusGetFromGroups)
		p.logger.Trace("开始获取源库Root Group")
		fromGroups, _, err := p.FromClient.Groups.ListGroups(&gitlab.ListGroupsOptions{
			TopLevelOnly: &topLevelOnly,
			ListOptions:  gitlab.ListOptions{PerPage: 100},
		})
		if err != nil {
			p.UpdateStatus(SyncPipeStatusGetFromGroupsError)
			return err
		}
		p.logger.WithField("FromGroups", fromGroups).Trace("完成获取源库Root Group")
		p.UpdateStatus(SyncPipeStatusGetToGroups)
		p.logger.Trace("开始获取目标库Root Group")
		toGroups, _, err := p.ToClient.Groups.ListGroups(&gitlab.ListGroupsOptions{
			TopLevelOnly: &topLevelOnly,
			ListOptions:  gitlab.ListOptions{PerPage: 100},
		})
		if err != nil {
			p.UpdateStatus(SyncPipeStatusGetToGroupsError)
			return err
		}
		p.logger.WithField("ToGroups", toGroups).Trace("完成获取目标库Root Group")
		p.SyncGroupList(fromGroups, toGroups, nil)
		p.logger.Debug("完成同步Root Group")
	} else {
		p.logger.
			WithField("ParentFrom", p.ShortGroup(parentFrom)).
			WithField("ParentTo", p.ShortGroup(parentTo)).
			Debug("开始同步Group")
		fromGroups, _, err := p.FromClient.Groups.ListSubGroups(parentFrom.ID, &gitlab.ListSubGroupsOptions{
			ListOptions: gitlab.ListOptions{PerPage: 100},
		})
		if err != nil {
			return err
		}
		p.logger.
			WithField("ParentFrom", p.ShortGroup(parentFrom)).
			WithField("FromGroups", p.ShortGroups(fromGroups)).
			Trace("获取到源库的Group信息")
		toGroups, _, err := p.ToClient.Groups.ListSubGroups(parentTo.ID, &gitlab.ListSubGroupsOptions{
			ListOptions: gitlab.ListOptions{PerPage: 100},
		})
		if err != nil {
			return err
		}
		p.logger.
			WithField("ParentTo", p.ShortGroup(parentTo)).
			WithField("ToGroups", p.ShortGroups(toGroups)).
			Trace("获取到目标库的Group信息")
		p.SyncGroupList(fromGroups, toGroups, parentTo)
	}
	return nil
}

func (p *SyncPipe) SyncGroupProject(formGroup *gitlab.Group, toGroup *gitlab.Group) error {
	p.logger.
		WithField("FromGroup", p.ShortGroup(formGroup)).
		WithField("ToGroup", p.ShortGroup(toGroup)).
		Debug("开始同步Group的Projects")
	p.UpdateStatus(SyncPipeStatusGetFromGroupProjects)
	fromProjects, _, err := p.FromClient.Groups.ListGroupProjects(formGroup.ID, &gitlab.ListGroupProjectsOptions{
		ListOptions: gitlab.ListOptions{PerPage: 100},
	})
	if err != nil {
		p.UpdateStatus(SyncPipeStatusGetFromGroupProjectsErr)
		return err
	}
	p.logger.
		WithField("FromGroup", p.ShortGroup(formGroup)).
		WithField("FromProjects", p.ShortProjects(fromProjects)).
		Trace("获取到源库Projects")
	p.UpdateStatus(SyncPipeStatusGetToGroupProjects)
	toProjects, _, err := p.ToClient.Groups.ListGroupProjects(toGroup.ID, &gitlab.ListGroupProjectsOptions{
		ListOptions: gitlab.ListOptions{PerPage: 100},
	})
	if err != nil {
		p.UpdateStatus(SyncPipeStatusGetFromGroupProjectsErr)
		return err
	}
	p.logger.
		WithField("ToGroup", p.ShortGroup(toGroup)).
		WithField("ToProjects", p.ShortProjects(toProjects)).
		Trace("获取到目标库Projects")
	for _, fromProject := range fromProjects {
		toProjectT := linq.From(toProjects).WhereT(
			func(p *gitlab.Project) bool {
				return p.Path == fromProject.Path
			}).First()
		var toProject *gitlab.Project
		if toProjectT != nil {
			toProject = toProjectT.(*gitlab.Project)
		} else {
			p.logger.
				WithField("FromGroup", p.ShortGroup(formGroup)).
				WithField("ToGroup", p.ShortGroup(toGroup)).
				WithField("FromProject", p.ShortProject(fromProject)).
				Debug("目标库的的Project不存在，开始创建")
			p.UpdateStatus(SyncPipeStatusCreateToProject)
			toProject, _, err = p.ToClient.Projects.CreateProject(&gitlab.CreateProjectOptions{
				Path: &fromProject.Path,
				Name: &fromProject.Name, NamespaceID: &toGroup.ID})
			if err != nil {
				p.UpdateStatus(SyncPipeStatusCreateToProjectErr)
				p.logger.WithField("NamespaceID", toGroup.ID).
					WithField("Path", fromProject.Path).
					WithField("Name", fromProject.Name).
					Error("CreateProject:", err)
				continue
			}
			p.logger.
				WithField("FromGroup", p.ShortGroup(formGroup)).
				WithField("ToGroup", p.ShortGroup(toGroup)).
				WithField("FromProject", p.ShortProject(fromProject)).
				Debug("完成目标库的的Project创建")
		}
		err := p.SyncProjectData(fromProject, toProject)
		if err != nil {
			p.logger.WithField("FromProject", p.ShortProject(fromProject)).
				WithField("ToProject", p.ShortProject(toProject)).
				Error("SyncProjectData:", err)
		}
	}
	p.logger.
		WithField("FromGroup", p.ShortGroup(formGroup)).
		WithField("ToGroup", p.ShortGroup(toGroup)).
		Debug("完成同步Group的Projects")
	return nil
}

func (p *SyncPipe) SyncProjectData(from *gitlab.Project, to *gitlab.Project) error {
	p.logger.WithField("FromProject", p.ShortProject(from)).
		WithField("ToProject", p.ShortProject(to)).
		Debug("开始同步Project数据")
	err, syncMap := p.SyncBranch(from, to)
	if err != nil {
		return err
	}
	needSync := false
	for fromBranch, toBranch := range syncMap {
		p.logger.WithField("FromBranch", p.ShortBranch(fromBranch)).
			WithField("ToBranch", p.ShortBranch(toBranch)).
			Trace("对比Commit")
		if strings.ToUpper(fromBranch.Commit.ID) != strings.ToUpper(toBranch.Commit.ID) {
			p.logger.WithField("FromBranch", p.ShortBranch(fromBranch)).
				WithField("ToBranch", p.ShortBranch(toBranch)).
				Debug("Commit不一样，同步分支")
			err := p.PushTo(from, to, fromBranch)
			if err != nil {
				p.logger.WithField("FromProject", p.ShortProject(from)).
					WithField("ToProject", p.ShortProject(to)).
					Error("PushTo:", err)
				//os.Exit(9)
			}
		}
	}
	if needSync {

	}

	syncMap = nil
	return nil
}
func (p *SyncPipe) ensureDir(dirName string) error {
	err := os.RemoveAll(dirName)
	if err != nil {
		p.logger.Warning(err)
	}
	err = os.MkdirAll(dirName, 0775)
	if err == nil {
		return nil
	}
	if os.IsExist(err) {
		// check that the existing path is a directory
		info, err := os.Stat(dirName)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return errors.New("path exists but is not a directory")
		}
		return nil
	}
	return err
}

func (p *SyncPipe) PushTo(from *gitlab.Project, to *gitlab.Project, fromBranch *gitlab.Branch) error {

	p.logger.
		WithField("FromProject", p.ShortProject(from)).
		WithField("ToProject", p.ShortProject(to)).
		Debug("开始同步库")
	err := p.ensureDir(configs.TempDir)
	if err != nil {
		return err
	}
	exeCmd := func(basePath string) types.Option {
		return git.CmdExecutor(func(ctx context.Context, name string, debug bool, args ...string) (string, error) {
			if debug {
				p.logger.Println(name, strings.Join(args, " "))
			}
			cmd := exec.CommandContext(ctx, name, args...)
			cmd.Dir = basePath
			output, err := cmd.CombinedOutput()

			return string(output), err
		})
	}
	p.logger.Trace("创建git存储库")
	fromUri, err := url.Parse(from.HTTPURLToRepo)
	if err != nil {
		return errors.New(fmt.Sprintf("HTTPURLToRepo parse:%s", err))
	}
	fromUri.User = url.UserPassword("gl2gl_sync", p.FromAccessToken)
	cloneText, err := git.Clone(exeCmd(configs.TempDir), clone.Repository(fromUri.String()),
		git.Debug, clone.Branch(fromBranch.Name))
	p.logger.Logger.Trace("cloneText:", cloneText)
	if err != nil {
		return errors.New(fmt.Sprintf("Clone:%s", err))
	}

	remote, err := git.Remote(exeCmd(path.Join(configs.TempDir, from.Path)))
	log.Println(remote)
	if err != nil {
		return err
	}
	fetchText, err := git.Fetch(exeCmd(path.Join(configs.TempDir, from.Path)), fetch.Verbose, fetch.All, git.Debug)
	p.logger.Logger.Trace("fetchText:", fetchText)
	if err != nil {
		return errors.New(fmt.Sprintf("Fetch:%s", err))
	}
	uri, err := url.Parse(to.HTTPURLToRepo)
	if err != nil {
		return errors.New(fmt.Sprintf("HTTPURLToRepo parse:%s", err))
	}
	uri.User = url.UserPassword("gl2gl_sync", p.ToAccessToken)

	p.logger.Trace("开始推送远程")
	pushText, err := git.Push(exeCmd(path.Join(configs.TempDir, from.Path)), push.Remote(uri.String()), push.All, git.Debug)
	p.logger.Logger.Trace("pushText:", pushText)
	if err != nil {
		return errors.New(fmt.Sprintf("toRemote Push:%s", err))
	}
	p.logger.Trace("完成推送本地分支到远程分支")
	return nil
}

func (p *SyncPipe) SyncBranch(from *gitlab.Project, to *gitlab.Project) (error, map[*gitlab.Branch]*gitlab.Branch) {
	p.logger.WithField("FromProject", p.ShortProject(from)).
		WithField("ToProject", p.ShortProject(to)).
		Debug("开始同步Project的分支")
	p.UpdateStatus(SyncPipeStatusGetFromBranches)
	p.logger.WithField("FromProject", p.ShortProject(from)).
		Trace("开始获取源库Project的分支")
	_fromBranches, _, err := p.FromClient.Branches.ListBranches(from.ID, &gitlab.ListBranchesOptions{
		ListOptions: gitlab.ListOptions{PerPage: 100},
	})
	if err != nil {
		p.UpdateStatus(SyncPipeStatusGetFromBranchesError)
		return err, nil
	}
	p.logger.WithField("FromProject", p.ShortProject(from)).
		WithField("FromBranches", p.ShortBranches(_fromBranches)).
		Trace("完成获取源库Project的分支")

	p.logger.WithField("ToProject", p.ShortProject(to)).
		Trace("开始获取目标库Project的分支")
	p.UpdateStatus(SyncPipeStatusGetToBranches)
	toBranches, _, err := p.ToClient.Branches.ListBranches(to.ID, &gitlab.ListBranchesOptions{
		ListOptions: gitlab.ListOptions{PerPage: 100},
	})
	if err != nil {
		p.UpdateStatus(SyncPipeStatusGetToBranchesError)
		return err, nil
	}
	p.logger.
		WithField("ToProject", p.ShortProject(to)).
		WithField("ToBranches", p.ShortBranches(toBranches)).
		Trace("完成获取目标库Project的分支")
	branchMap := make(map[*gitlab.Branch]*gitlab.Branch)
	for _, fromBranch := range _fromBranches {
		toBranchT := linq.From(toBranches).WhereT(
			func(b *gitlab.Branch) bool {
				return b.Name == fromBranch.Name
			}).First()
		var toBranch *gitlab.Branch
		if toBranchT == nil {
			p.logger.
				WithField("FromProject", p.ShortProject(from)).
				WithField("ToProject", p.ShortProject(to)).
				WithField("FromBranch", p.ShortBranch(fromBranch)).
				Debug("目标库的的分支不存在，开始创建")
			err := p.PushTo(from, to, fromBranch)
			if err != nil {
				p.logger.WithField("FromProject", p.ShortProject(from)).
					WithField("ToProject", p.ShortProject(to)).
					WithField("FromBranch", p.ShortBranch(fromBranch)).
					Error("PushTo:", err)
				continue
			}
			toBranch, _, err = p.ToClient.Branches.GetBranch(to.ID, fromBranch.Name)
			if err != nil {
				p.logger.WithField("ToProject", p.ShortProject(to)).
					WithField("FromBranch", p.ShortBranch(fromBranch)).
					Error("GetBranch:", err)
				continue
			}
			p.logger.
				WithField("FromProject", p.ShortProject(from)).
				WithField("ToProject", p.ShortProject(to)).
				WithField("ToBranch", p.ShortBranch(toBranch)).
				Debug("完成目标库的的分支创建")

		} else {
			toBranch = toBranchT.(*gitlab.Branch)
		}
		branchMap[fromBranch] = toBranch
	}
	p.logger.WithField("FromProject", p.ShortProject(from)).
		WithField("ToProject", p.ShortProject(to)).
		Debug("完成始同步Project的分支")
	_fromBranches = nil
	return nil, branchMap
}

func (p *SyncPipe) SyncGroupList(fromGroups []*gitlab.Group, toGroups []*gitlab.Group, toParentGroup *gitlab.Group) {
	for _, fromGroup := range fromGroups {
		toGroupT := linq.From(toGroups).
			WhereT(func(g *gitlab.Group) bool {
				return strings.ToUpper(g.Path) == strings.ToUpper(fromGroup.Path)
			}).First()
		var err error
		var toGroup *gitlab.Group
		if toGroupT == nil {
			p.logger.WithField("FromGroup", p.ShortGroup(fromGroup)).Debug("目标库Group不存在,开始创建")
			p.UpdateStatus(SyncPipeStatusCreateToGroup)
			var parentId *int
			if toParentGroup != nil {
				parentId = &toParentGroup.ID
			}
			toGroup, _, err = p.ToClient.Groups.CreateGroup(&gitlab.CreateGroupOptions{
				Path: &fromGroup.Path, Name: &fromGroup.Name, ParentID: parentId,
			})
			if err != nil {
				p.UpdateStatus(SyncPipeStatusCreateToGroupErr)
				p.logger.WithField("FromGroup", p.ShortGroup(fromGroup)).
					Error("CreateGroup:", err)
				continue
			}
			p.logger.WithField("ToGroup", p.ShortGroup(toGroup)).Debug("完成目标库Group创建")
		} else {
			toGroup = toGroupT.(*gitlab.Group)
		}
		err = p.SyncGroupProject(fromGroup, toGroup)
		if err != nil {
			p.logger.
				WithField("FromGroup", p.ShortGroup(fromGroup)).
				WithField("ToGroup", p.ShortGroup(toGroup)).
				Error("SyncGroupProject:", err)
			continue
		}
		err = p.SyncGroup(fromGroup, toGroup)
		if err != nil {
			p.logger.
				WithField("FromGroup", p.ShortGroup(fromGroup)).
				WithField("ToGroup", p.ShortGroup(toGroup)).
				Error("SyncGroup:", err)
		}
	}
}

func (p *SyncPipe) SyncUser() error {
	p.logger.Debug("开始同步用户..")
	p.logger.Trace("开始获取源库的用户列表")
	fromUsers, _, err := p.FromClient.Users.ListUsers(&gitlab.ListUsersOptions{
		ListOptions: gitlab.ListOptions{PerPage: 100},
	})
	p.logger.Trace("获取到源库的用户列表:", p.ShortUsers(fromUsers))
	if err != nil {
		return err
	}
	p.logger.Trace("开始获取目标库的用户列表")
	toUsers, _, err := p.ToClient.Users.ListUsers(&gitlab.ListUsersOptions{
		ListOptions: gitlab.ListOptions{PerPage: 100},
	})
	p.logger.Trace("获取到目标库的用户列表:", p.ShortUsers(toUsers))
	if err != nil {
		return err
	}
	for _, fromUser := range fromUsers {
		toUserT := linq.From(toUsers).
			WhereT(func(g *gitlab.User) bool {
				return strings.ToUpper(g.Username) == strings.ToUpper(fromUser.Username)
			}).First()
		var toUser *gitlab.User
		if toUserT == nil {
			randomPassword := true
			p.logger.WithField("FromUser", p.ShortUser(fromUser)).Debug("目标库的的用户不存在，开始创建")
			toUser, _, err = p.ToClient.Users.CreateUser(&gitlab.CreateUserOptions{
				Email:               &fromUser.Email,
				Name:                &fromUser.Name,
				Username:            &fromUser.Username,
				ForceRandomPassword: &randomPassword,
			})
			if err != nil {
				p.logger.
					WithField("FromUser", p.ShortUser(fromUser)).
					Error("CreateUser:", err)
				continue
			}
			p.logger.WithField("ToUser", p.ShortUser(toUser)).Debug("完成目标库的用户创建")

		} else {
			toUser = toUserT.(*gitlab.User)
		}
		err := p.SyncUserProjects(fromUser, toUser)
		if err != nil {
			p.logger.
				WithField("FromUser", p.ShortUser(fromUser)).
				WithField("ToUser", p.ShortUser(toUser)).
				Error("SyncUserProjects:", err)
		}

	}
	return nil
}

func (p *SyncPipe) SyncUserProjects(fromUser *gitlab.User, toUser *gitlab.User) error {
	p.logger.
		WithField("FromUser", p.ShortUser(fromUser)).
		WithField("ToUser", p.ShortUser(toUser)).
		Debug("开始同步用户Projects")
	p.logger.
		WithField("FromUser", p.ShortUser(fromUser)).
		Trace("获取源库用户的Projects")
	formProjects, _, err := p.FromClient.Projects.ListUserProjects(fromUser.ID,
		&gitlab.ListProjectsOptions{ListOptions: gitlab.ListOptions{PerPage: 100}},
	)
	if err != nil {
		return err
	}
	p.logger.WithField("FromUser", p.ShortUser(fromUser)).
		WithField("FormProjects", p.ShortProjects(formProjects)).
		Trace("完成源库用户的Projects获取")
	p.logger.WithField("ToUser", p.ShortUser(toUser)).Trace("获取目标库用户的Projects")
	toProjects, _, err := p.ToClient.Projects.ListUserProjects(toUser.ID,
		&gitlab.ListProjectsOptions{ListOptions: gitlab.ListOptions{PerPage: 100}},
	)
	if err != nil {
		return err
	}
	p.logger.WithField("ToUser", p.ShortUser(toUser)).
		WithField("ToProjects", p.ShortProjects(toProjects)).
		Trace("完成目标库用户的Projects获取")
	for _, fromProject := range formProjects {
		toProjectT := linq.From(toProjects).
			WhereT(func(g *gitlab.Project) bool {
				return strings.ToUpper(g.PathWithNamespace) == strings.ToUpper(fromProject.PathWithNamespace)
			}).First()
		var toProject *gitlab.Project
		if toProjectT == nil {
			p.logger.WithField("ToUser", p.ShortUser(toUser)).
				WithField("FromProject", p.ShortProject(fromProject)).
				Debug("目标库的用户Project不存在，开始创建")
			toProject, _, err = p.ToClient.Projects.CreateProjectForUser(
				toUser.ID,
				&gitlab.CreateProjectForUserOptions{
					Name: &fromProject.Name,
					Path: &fromProject.Path,
				})
			p.logger.WithField("ToUser", p.ShortUser(toUser)).
				WithField("ToProject", p.ShortProject(toProject)).
				Debug("完成目标库的用户Project创建")
			if err != nil {
				p.logger.
					WithField("ToUser", p.ShortUser(toUser)).
					WithField("FromProject", p.ShortProject(fromProject)).
					Error("CreateProjectForUser:", err)
				continue
			}
		} else {
			toProject = toProjectT.(*gitlab.Project)
		}
		err := p.SyncProjectData(fromProject, toProject)
		if err != nil {
			p.logger.
				WithField("FromProject", p.ShortProject(fromProject)).
				WithField("ToProject", p.ShortProject(toProject)).
				Error("SyncProjectData:", err)
			continue
		}
	}
	p.logger.
		WithField("FromUser", p.ShortUser(fromUser)).
		WithField("ToUser", p.ShortUser(toUser)).
		Debug("完成同步用户Projects")
	return nil
}

func (p *SyncPipe) ShortGroup(group *gitlab.Group) string {
	return fmt.Sprintf("ParentId:%d,GroupId:%d,GroupName:%s,WebURL:%s,FullPath:%s",
		group.ParentID,
		group.ID,
		group.Name,
		group.WebURL,
		group.FullPath)
}

func (p *SyncPipe) ShortProjects(projects []*gitlab.Project) string {
	var buffer bytes.Buffer
	for i, project := range projects {
		buffer.WriteString("{")
		buffer.WriteString(p.ShortProject(project))
		buffer.WriteString("}")
		if i < len(projects)-1 {
			buffer.WriteString(",")
		}
	}
	return buffer.String()
}

func (p *SyncPipe) ShortProject(project *gitlab.Project) string {
	return fmt.Sprintf("NamespaceId:%d,ID:%d,Name:%s,WebURL:%s,PathWithNamespace:%s",
		project.Namespace.ID,
		project.ID,
		project.Name,
		project.WebURL,
		project.PathWithNamespace)
}

func (p *SyncPipe) ShortBranch(branch *gitlab.Branch) string {
	return fmt.Sprintf("Name:%s,WebURL:%s,Commit.ID:%s,Commit.Messag:%s",
		branch.Name,
		branch.WebURL,
		branch.Commit.ID,
		branch.Commit.Message)
}

func (p *SyncPipe) ShortBranches(branches []*gitlab.Branch) string {
	var buffer bytes.Buffer
	for i, branch := range branches {
		buffer.WriteString("{")
		buffer.WriteString(p.ShortBranch(branch))
		buffer.WriteString("}")
		if i < len(branches)-1 {
			buffer.WriteString(",")
		}
	}
	return buffer.String()
}

func (p *SyncPipe) ShortUsers(users []*gitlab.User) string {
	var buffer bytes.Buffer
	for i, user := range users {
		buffer.WriteString("{")
		buffer.WriteString(p.ShortUser(user))
		buffer.WriteString("}")
		if i < len(users)-1 {
			buffer.WriteString(",")
		}
	}
	return buffer.String()
}

func (p *SyncPipe) ShortUser(users *gitlab.User) string {
	return fmt.Sprintf("Username:%s,Email:%s,Name:%s,WebURL:%s",
		users.Username,
		users.Email,
		users.Name,
		users.WebURL)
}

func (p *SyncPipe) ShortGroups(groups []*gitlab.Group) interface{} {
	var buffer bytes.Buffer
	for i, group := range groups {
		buffer.WriteString("{")
		buffer.WriteString(p.ShortGroup(group))
		buffer.WriteString("}")
		if i < len(groups)-1 {
			buffer.WriteString(",")
		}
	}
	return buffer.String()
}

//func (p *SyncPipe) makeRemoteCallbacks(name, accessToken string) *git.RemoteCallbacks {
//	return &git.RemoteCallbacks{
//		SidebandProgressCallback: func(str string) error {
//			p.logger.Logger.Trace(fmt.Sprintf("[%s]===>", name), "SidebandProgressCallback", str)
//			return nil
//		},
//		CompletionCallback: func(c git.RemoteCompletion) error {
//			p.logger.Logger.Trace(fmt.Sprintf("[%s]===>", name), "CompletionCallback", c)
//			return nil
//		},
//		CredentialsCallback: func(url string, username_from_url string, allowed_types git.CredentialType) (*git.Credential, error) {
//			cred, err := git.NewCredentialUserpassPlaintext("gl2gl_sync", accessToken)
//			errStr := ""
//			if err != nil {
//				errStr = err.Error()
//			}
//			p.logger.Logger.Trace(fmt.Sprintf("[%s]===>", name), fmt.Sprintf("CredentialsCallback url:%s,username_from_url:%s,err:%s", url, username_from_url, errStr))
//			return cred, err
//		},
//		TransferProgressCallback: func(stats git.TransferProgress) error {
//			p.logger.Logger.Trace(fmt.Sprintf("[%s]===>", name), fmt.Sprintf(
//				"TransferProgressCallback IndexedObjects:%d,LocalObjects:%d,TotalObjects:%d,TotalDeltas:%d,ReceivedObjects:%d,ReceivedBytes:%d",
//				stats.IndexedObjects, stats.LocalObjects, stats.TotalObjects, stats.TotalDeltas, stats.ReceivedObjects, stats.ReceivedBytes,
//			))
//			return nil
//		},
//		UpdateTipsCallback: func(refname string, a *git.Oid, b *git.Oid) error {
//			p.logger.Logger.Trace(fmt.Sprintf("[%s]===>", name), fmt.Sprintf("CompletionCallback refname:%s,a:%s,b:%s", refname, a.String(), b.String()))
//			return nil
//		},
//		CertificateCheckCallback: func(cert *git.Certificate, valid bool, hostname string) error {
//			p.logger.Logger.Trace(fmt.Sprintf("[%s]===>", name), fmt.Sprintf("CertificateCheckCallback cert:%v,CertificateCheckCallback:%s", cert.X509, hostname))
//			return nil
//		},
//		PackProgressCallback: func(stage int32, current, total uint32) error {
//			p.logger.Logger.Trace(fmt.Sprintf("[%s]===>", name), fmt.Sprintf("PackProgressCallback stage:%d,current:%d,total:%d", stage, current, total))
//			return nil
//		},
//		PushTransferProgressCallback: func(current, total uint32, bytes uint) error {
//			p.logger.Logger.Trace(fmt.Sprintf("[%s]===>", name), fmt.Sprintf("PushTransferProgressCallback current:%d,total:%d,bytes:%d", current, total, bytes))
//			return nil
//		},
//		PushUpdateReferenceCallback: func(refname, status string) error {
//			p.logger.Logger.Trace(fmt.Sprintf("[%s]===>", name), fmt.Sprintf("PushUpdateReferenceCallback refname:%s,status:%s", refname, status))
//			return nil
//		},
//	}
//}
