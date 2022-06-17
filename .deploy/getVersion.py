import re,sys,datetime







verFile = open(sys.argv[1], "r+")
text=verFile.read()
now =  datetime.datetime.now()
nowStr=now.strftime("%Y-%m-%d %H:%M:%S")
reg=re.compile('(?<=VERSION)[^\"]+\"[^\"]+\"')
mg=re.search(reg,text)
print(mg.group(0).replace('=','').replace('\"','').strip())
text = re.sub(r'PublishTime[^\"]+\"[^\"]+\"', 'PublishTime string = "'+nowStr+'"', text, flags=re.DOTALL)
# print(type(text))
# print(text)
verFile.seek(0)
verFile.write(text)
verFile.close()