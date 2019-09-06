# brew
```
brew list 							(查看已经安装的包)
brew update                         (更新 HomeBrew)
brew outdated 						(查看可更新的包)
brew upgrade                        (更新所有包)
brew upgrade formula_name 			(更新指定包)
brew cleanup                        (清理所有包的旧版本)
brew cleanup formula_name 			(清理指定包的旧版本)
brew cleanup -n                     (查看可清理的旧版本包，不执行实际操作)
brew uninstall formula_name --force	(彻底卸载指定软件包)
brew pin formula_name 				(锁定指定包)
brew unpin formula_name             (取消锁定)
brew info 							(显示安装包信息)
brew info formula_name 				(显示指定包的信息)
brew deps formula_name 				(显示指定包依赖关系)
brew deps --installed --tree        (查看已安装的包的依赖，树形显示)
```