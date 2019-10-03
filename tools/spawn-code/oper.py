import os 

class directory:
        def name(self):
            direct = os.getcwd()
            idx = direct.rfind("\\")
            if idx <= 0 :
                idx = direct.rfind("/")
            if idx <= 0:
                return direct
            return direct[idx + 1:]


class oper_spawn:
    def __init__(self, filename, pkgname):
        self.fullpath = filename
        self.pkgname = pkgname
    def read(self):
        try:
            f = open(self.fullpath, "r")
            result = f.read()
            f.close()
            return result
        except FileNotFoundError:
            return ""
        except IOError:
            return ""

    def write(self, data):
        try:
            f = open(self.fullpath, "w")
            f.write(data)
            f.flush()
            f.close()
        except IOError:
            print("File is not accessible.")
            exit(0)


class read_template:

    def __init__(self, filename):
        self.fullpath = filename

    def read(self):
        try:
            f = open(self.fullpath, "r")
            result = f.read()
            f.close()
            return result
        except FileNotFoundError:
            print(self.fullpath + " File s not found.")
            exit(0)
        except IOError:
            print("File is not accessible.")
            exit(0)

    def sreach_func(self, desc):
        start = 0
        result = []      
        while(True):
            nstart = desc.find("/[", start)
            
            if (nstart <= 0):
                break
            nend = desc.find("]/", start)
            if (nend <= 0):
                break 
            start = nend + 2
            result.append(desc[nstart + 2: nend])
        return result

    def replace_name(self, name, desc):
        return desc.replace("[label_struct]", name)

    def replace_note(self, note, desc):
        return desc.replace("[label_note]", note)


class read_spawn:

    def __init__(self, filename):
        self.fullpath = filename

    def read(self):
        try:
            f = open(self.fullpath, "r")
            result = f.readlines()
            f.close()
            return result
        except FileNotFoundError:
            print(self.fullpath + " File s not found.")
            exit(0)
        except IOError:
            print("File is not accessible.")
            exit(0)
    