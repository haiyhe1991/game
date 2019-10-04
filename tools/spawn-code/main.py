import sys
import os
import oper
import datetime

#[project name],[function name],[struct name],[note]
welcome = "Server System Logic Code Template Generator v1.0.0"
symbold = "\\"
paths   = symbold + "tools" + symbold

if "__main__" == __name__:  
    workpath = os.getcwd()
    package  = oper.directory().name()
    configfile   = workpath + paths + "config.txt"
    templatefile = workpath + paths + "template.txt"
    spawnfile    = workpath + paths + "spawn.txt"
    inittemplatefile = workpath + paths + "init-template.txt"


    ifile = 0

    print (welcome)
    print ("work directory:" + workpath)
    print ("package name:" + package)
    print ("read template file:" + templatefile)
    template_data = oper.read_template(templatefile).read()
    print ("read init template file:" + inittemplatefile)
    init_template_data =  oper.read_template(inittemplatefile).read()
    print ("read spawn file:" + spawnfile)
    spawn_data    = oper.read_spawn(spawnfile).read()
    print ("read config file:" + configfile)
    config_data   = oper.read_spawn(configfile).read()
    
    #生成INIT文件
    spawn_new_data  = ""
    spawn_init_file = workpath + symbold + config_data[0] + "-init.go"
    for row in range(len(spawn_data)):
        info = spawn_data[row]
        key = info.strip().split(",")
        if len(key) != 4:
            continue
        filename = key[0] + "-" + key[1] + "-process.go"
        filepath = workpath + symbold + filename
        print("spawn file:" + filepath)
        filehandle = oper.oper_spawn(filepath, package)
        sur = filehandle.read()
        new_template_data = oper.read_template(templatefile).replace_name(key[2], template_data)
        new_template_data = oper.read_template(templatefile).replace_note(key[3], new_template_data)
        if (ifile> 0):
            spawn_new_data += "\n"
        spawn_new_data +=  "	module.FactoryInstance().Register(\"" + package + "." + key[2] + "\",&" + key[2] + "{})"
       
        if sur == "":
             sur += "package " + package + "\n\n"
             sur += new_template_data
  
        oper.oper_spawn(filepath, package).write(sur)
        ifile += 1
        print ("spawn file:" + filepath + " complate")
    init_template_data =  "package " + package + "\n\n" +  init_template_data.replace("[label list]", spawn_new_data)
    oper.oper_spawn(spawn_init_file, package).write(init_template_data)
    print("files:" + str(ifile))


    