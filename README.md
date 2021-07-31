## 学习 interpreter 的笔记
跟着 `Writing An Interpreter In Go` 这本书进行学习


### 知识点

#### parser generator
常见的解释生成器工具有: `yacc`, `bison`, `antler`  
大多情况下直接使用现有的`parser generator`就可以了不用重新造轮子

#### context-free grammar (CFG)
CFG是一组规则,用于描述一个程序语言的规则(语法)  
常用的描述有: `Backus-Naur Form (BNF)`, `Extended Backus-Naur Form (EBNF)`

#### 解析器
解析程序语言有两种主要的策略 
1. `top-down parsing` strategies
2. `bottom upparsing` strategies
本书的解释器采用 `top-down`, 叫做 `pratt parser` 

主要目的先构建AST(abstract syntax tree)

```
        Node
Statement Expression
```

AST上的每个节点都是`Node`接口的实现, 都有`TokenLiteral()`


#### 测试工具的使用

##### GoMock
项目地址: <https://github.com/golang/mock>
教程地址 <https://www.jianshu.com/p/f4e773a1b11f> 
1 安装
```$xslt
go get github.com/golang/mock/gomock
go get github.com/golang/mock/mockgen
```

2 定义接口,生成对应mock文件
```$xslt
# 使用mockgen生成mock文件 (Source mode)
mockgen -source=lexer/lexer.go --destination=lexer/mocks/mock_lexer.go -package=mocks
```

3 测试
```
GO_MOCK_TEST=1 go test -v
```

##### GoConvey
项目地址: <https://github.com/smartystreets/goconvey>
管理和运行测试用例, 提供丰富的断言函数, 还提供webui的测试框架

运行 `$GOPATH/bin/goconvey` 可以打开web界面

##### GoStub
项目地址: <https://github.com/prashantv/gostub>
<https://www.jianshu.com/p/53a531852619>
可以对全局变量、函数或过程打桩