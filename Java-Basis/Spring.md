# IOC

## 简约步骤

### 按照Spring IOC的加载流程

1. 找到bean
>找到bean在什么地方，是对BeanDefinition的资源定位，是由ResourceLoader通过统一的Resource接口来完成，这个接口对各中形式的Resource都提供了统一接口，比如Xml，比如annotation。而这些都是由ResourceLoader来完成的

2. 载入并注册bean
>找到bean后，将bean注册到我们的IOC容器中。Spring是通过一些ApplicationContext来完成的，比如FileSystemXmlApplicationContext, ClassPathXmlApplicationContext以及我们最常见的XmlWebApplicationContext，读取之后将bean注册到IOC容器中，简单来说，就是把读取的bean都放到一个map中。

3. 注入bean
>当我们要用bean时，由IOC容器自动的注入进去。