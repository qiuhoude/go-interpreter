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
AST上的每个节点都是`Node`接口的实现, 都有`TokenLiteral()`  
```
        Node
Statement Expression
```


`Top Down Operator Precedence`  在1973年发布,作者:`Vaughan Pratt`   
`Game Programming Patterns` 被书作者强烈推荐

##### 术语
1. prefix operator (前缀操作), `--foobar`
2. postfix operator(后缀操作), `foobar++`
3. infix operators(中缀操作), `5 + 5 * 10`, 有优先级,需要定义 precedences

主要思路, 根据token type关联一些解析函数, 函数返回都是 AST node expression,
每次遇到对应的token type都会找出对应的解析函数, 
每个token type最多可以有两个解析函数相关联,具体取决于token的位置是prefix 或 infix


#### Evaluation
eval 会有不同的 strategies
遍历AST树叫做 `tree-walking interpreters` (最慢的方式)  
会对AST进行重写(删除无用的结构)来达到一些小优化  
或者转换成其他intermediate representation(IR) 中间表现层  
  
解释器递归AST多次进入某个分支会将该分支编译成机器码

NULL 的引入称之为亿万美元错误 `billion-dollar mistake` (billion十亿,有连接符号dollar不用加s)   

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