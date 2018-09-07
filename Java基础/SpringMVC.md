# Spring MVC

## Spring MVC 框架有什么用？优点？
Spring Web MVC 框架提供 **模型**-**视图**-**控制器** 架构和随时可用的组件，用于开发 **灵活且松散耦合**的web应用程序。 输入逻辑、业务逻辑、UI逻辑。

## DispatcherServlet

![DispatcherServlet的工作流程](Images/SpringMVC之DispatcherServlet的工作流程.jpg)

1. 向服务器发送 HTTP 请求，请求被前端控制器 DispatcherServlet 捕获。

2. DispatcherServlet 根据 `servlet.xml` 中的配置对请求的 URL 进行解析——>得到请求资源标识符（URI）——>根据该 URI，调用 `HandlerMapping` 获得该 `Handler` 配置的所有相关的对象（包括 Handler 对象以及 Handler 对象对应的拦截器）——>以`HandlerExecutionChain` 对象的形式返回。

3. DispatcherServlet 根据获得的Handler，选择一个合适的`HandlerAdapter`.
   >附注：如果成功获得HandlerAdapter后，此时将开始执行拦截器的 `preHandler(...)`方法）。

4. 提取`Request`中的模型数据，填充`Handler入参`，开始执行`Handler（Controller)`。 在填充Handler的入参过程中，根据你的配置，Spring 将帮你做一些额外的工作：

   * HttpMessageConveter： 将请求消息（如 Json、xml 等数据）转换成一个对象，将对象转换为指定的响应信息。

    * 数据转换：对请求消息进行数据转换。如String转换成Integer、Double等。

    * 数据根式化：对请求消息进行数据格式化。 如将字符串转换成格式化数字或格式化日期等。

    * 数据验证： 验证数据的有效性（长度、格式等），验证结果存储到BindingResult或Error中。

5. Handler(Controller)执行完成后，向 DispatcherServlet 返回一个`ModelAndView` 对象；

6. 根据返回的ModelAndView，选择一个适合的 `ViewResolver`（必须是已经注册到 Spring 容器中的ViewResolver)返回给DispatcherServlet。

7. `ViewResolver` 结合Model和View，来**渲染视图**。

8. 视图负责将渲染结果返回给客户端。