package services

import (
	"context"
	"database/sql"
	"drtech.co/gl2gl/orm"
	"fmt"
	"github.com/ahmetb/go-linq/v3"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/xanzy/go-gitlab"
	"strings"
	"time"
)

type SyncPipe struct {
	ConfigId        int64
	FromAddress     string
	FromAccessToken string
	FromClient      *gitlab.Client
	ToAddress       string
	ToAccessToken   string
	ToClient        *gitlab.Client
	Status          SyncPipeStatus

	groupMap   map[*gitlab.Group]*gitlab.Group
	projectMap map[*gitlab.Project]*gitlab.Project
	userMap    map[*gitlab.User]*gitlab.User
}

func (p *SyncPipe) Stop() error {
	return nil
}

func (p *SyncPipe) Run() error {
	p.groupMap = make(map[*gitlab.Group]*gitlab.Group)
	p.projectMap = make(map[*gitlab.Project]*gitlab.Project)
	p.userMap = make(map[*gitlab.User]*gitlab.User)
	p.UpdateStatus(SyncPipeStatusIniting)
	fromClient, err := gitlab.NewClient(p.FromAccessToken, gitlab.WithBaseURL(
		fmt.Sprintf("%s/api/v4", p.FromAddress),
	))
	if err != nil {
		//TODO log
		p.UpdateStatus(SyncPipeStatusFromClientInitError)
		return err
	}
	toClient, err := gitlab.NewClient(p.ToAccessToken, gitlab.WithBaseURL(
		fmt.Sprintf("%s/api/v4", p.ToAddress),
	))
	if err != nil {
		//TODO log
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
		Where(fromToConfigM.ID.Eq(sql.NullInt64{Int64: p.ConfigId, Valid: true})).Update(fromToConfigM.Status, status)
	if err != nil {
		//TODO log
	}
}

func (p *SyncPipe) Sync() {
	for {
		p.SyncGroup(nil)
		p.UpdateStatus(SyncPipeStatusGetFromProjects)
		fromProjects, _, err := p.FromClient.Projects.ListProjects(&gitlab.ListProjectsOptions{})
		if err != nil {
			//TODO log
			p.UpdateStatus(SyncPipeStatusGetFromProjectsError)
			time.Sleep(10 * time.Second)
			break
		}
		for _, project := range fromProjects {
			fmt.Println(project.PathWithNamespace)
		}
		fmt.Println(fromProjects)
		//toProjects, _, err := p.ToClient.Projects.ListProjects(&gitlab.ListProjectsOptions{})
		//if err != nil {
		//	//TODO log
		//	p.UpdateStatus(SyncPipeStatusGetToProjectsError)
		//	time.Sleep(10 * time.Second)
		//	break
		//}
		//fromProjectMap=make()
		//for _, project := range fromProjects {
		//	project.PathWithNamespace
		//}
	}
}

func (p *SyncPipe) SyncGroup(parent *gitlab.Group) error {

	if parent == nil {
		topLevelOnly := true
		p.UpdateStatus(SyncPipeStatusGetFromGroups)
		fromGroups, _, err := p.FromClient.Groups.ListGroups(&gitlab.ListGroupsOptions{
			TopLevelOnly: &topLevelOnly,
			ListOptions:  gitlab.ListOptions{PerPage: 100},
		})
		if err != nil {
			//TODO log
			p.UpdateStatus(SyncPipeStatusGetFromGroupsError)
			time.Sleep(10 * time.Second)
		}
		p.UpdateStatus(SyncPipeStatusGetToGroups)
		toGroups, _, err := p.ToClient.Groups.ListGroups(&gitlab.ListGroupsOptions{
			TopLevelOnly: &topLevelOnly,
			ListOptions:  gitlab.ListOptions{PerPage: 100},
		})
		if err != nil {
			//TODO log
			p.UpdateStatus(SyncPipeStatusGetToGroupsError)
			time.Sleep(10 * time.Second)
		}
		for _, fromGroup := range fromGroups {
			toGroupT := linq.From(toGroups).
				WhereT(func(g *gitlab.Group) bool {
					return strings.ToUpper(g.Path) == strings.ToUpper(fromGroup.Path)
				}).First()
			var toGroup *gitlab.Group
			if toGroup == nil {
				p.UpdateStatus(SyncPipeStatusCreateToGroup)
				toGroup, _, err = p.ToClient.Groups.CreateGroup(&gitlab.CreateGroupOptions{Path: &fromGroup.Path, Name: &fromGroup.Name})
				if err != nil {
					p.UpdateStatus(SyncPipeStatusCreateToGroupErr)
					//TODO log
					continue
				}
			} else {
				toGroup = toGroupT.(*gitlab.Group)
			}
			p.groupMap[fromGroup] = toGroup
			p.SyncGroupProject(fromGroup, toGroup)
		}

	} else {

	}

}

func (p *SyncPipe) SyncGroupProject(formGroup *gitlab.Group, toGroup *gitlab.Group) error {
	p.UpdateStatus(SyncPipeStatusGetFromGroupProjects)
	fromProjects, _, err := p.FromClient.Groups.ListGroupProjects(formGroup.ID, &gitlab.ListGroupProjectsOptions{
		ListOptions: gitlab.ListOptions{PerPage: 100},
	})
	if err != nil {
		p.UpdateStatus(SyncPipeStatusGetFromGroupProjectsErr)
		return err
	}
	p.UpdateStatus(SyncPipeStatusGetToGroupProjects)
	toProjects, _, err := p.ToClient.Groups.ListGroupProjects(toGroup.ID, &gitlab.ListGroupProjectsOptions{
		ListOptions: gitlab.ListOptions{PerPage: 100},
	})
	if err != nil {
		p.UpdateStatus(SyncPipeStatusGetFromGroupProjectsErr)
		return err
	}
	for _, fromProject := range fromProjects {
		toProjectT := linq.From(toProjects).WhereT(
			func(p *gitlab.Project) bool {
				return p.Path == fromProject.Path
			}).First()
		var toProject *gitlab.Project
		if toProjectT != nil {
			toProject = toProjectT.(*gitlab.Project)
		} else {
			p.UpdateStatus(SyncPipeStatusCreateToProject)
			toProject, _, err = p.ToClient.Projects.CreateProject(&gitlab.CreateProjectOptions{
				Path: &fromProject.Path,
				Name: &fromProject.Name})
			if err != nil {
				p.UpdateStatus(SyncPipeStatusCreateToProjectErr)
				continue
				//TODO log
			}
		}
		p.projectMap[toProject] = toProject
		p.SyncProjectData(fromProject, toProject)
	}
}

func (p *SyncPipe) SyncProjectData(from *gitlab.Project, to *gitlab.Project) error {
	p.UpdateStatus(SyncPipeStatusGetFromBranches)
	fromBranches, _, err := p.FromClient.Branches.ListBranches(from.ID, &gitlab.ListBranchesOptions{
		ListOptions: gitlab.ListOptions{PerPage: 100},
	})
	if err != nil {
		p.UpdateStatus(SyncPipeStatusGetFromBranchesError)
		return err
	}
	p.UpdateStatus(SyncPipeStatusGetToBranches)
	toBranches, _, err := p.ToClient.Branches.ListBranches(to.ID, &gitlab.ListBranchesOptions{
		ListOptions: gitlab.ListOptions{PerPage: 100},
	})
	if err != nil {
		p.UpdateStatus(SyncPipeStatusGetToBranchesError)
		return err
	}
	for _, fromBranch := range fromBranches {
		toBranchT := linq.From(toBranches).WhereT(
			func(b *gitlab.Branch) bool {
				return b.Name == fromBranch.Name
			}).First()
		var toBranch *gitlab.Branch
		if toBranchT != nil {
			toBranch = toBranchT.(*gitlab.Branch)
		} else {
			repo, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
				Auth: &http.BasicAuth{
					Username: "none", // yes, this can be anything except an empty string
					Password: p.FromAccessToken,
				},
				SingleBranch: true,
				URL:          from.HTTPURLToRepo,
				RemoteName:   "from",
			})
			_, err = repo.Branch(fromBranch.Name)
			if err != nil {
				continue
				//TODO log 不存在
			}
			refToPush := "tagName1^{}"
			refspecs := []config.RefSpec{config.RefSpec(refToPush + ":refs/heads/master")}

			err = repo.Push(&git.PushOptions{
				RemoteName: "origin",
				RefSpecs:   refspecs,
			})
			repo.Push(&git.PushOptions{
				RemoteName:        "to",
				RefSpecs:          nil,
				Auth:              nil,
				Progress:          nil,
				Prune:             false,
				Force:             false,
				InsecureSkipTLS:   false,
				CABundle:          nil,
				RequireRemoteRefs: nil,
			})
		}

	}
}
