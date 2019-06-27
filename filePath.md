1. 当前文件包名
	dir, _ := os.Getwd()
	fmt.Println(dir)

2. 当前文件全路径
	_, file, _, _ := runtime.Caller(0)
	fmt.Println(file)

3. 根据相对路径获取绝对路径
	absPath, _ := filepath.Abs(path)
	fmt.Println(absPath)

4. 获取文件名/包名
	name := filepath.Base(file/dir)
	fmt.Println(name)

5. 根据文件路径获取文件名
	fileName := filepath.Base(file)
	fmt.Println(fileName)

6. 根据文件路径获取包路径
	packageName := filepath.Dir(file)
	fmt.Println(packageName)

