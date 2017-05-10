# breeze-core

## 介绍

*breeze_core* 是 *breeze* 量化平台的核心组件，它主要负责代码管理、策略组装、上下文和事件驱动管理。 *breeze_core* 中涉及的主要概念有：

1. *策略*，策略的核心是一张 *事件* -> ActionList 的映射表，该表告诉平台，在某个特定的事件发生时，应该执行什么操作。
2. *UserContext*（用户上下文），一个 *UserContext* 业务上常常叫做一个 *组合*， *UserContext* 是某个特定 *策略* 的一个实例，在物理层面，一个 *UserContext* 就是一个 *策略* 加上 一组 *配置*。 *UserContext* 负责响应事件，它可以看做一个 *Actor*， 在我们的实现中，一个 *UserContext* 对应一个goroutine 运行的 *事件循环（event loop）*。
3. *事件*，很容易理解，不多赘言。
4. *Action*， *Action* 是一个类型为 `Event -> List<Event>` 的 *方法*， *Action* 这个定义是语言无关的， *Action* 甚至不必对应一段代码，只要满足类型约束即可。
5. *ActionList*,  *Action* 的列表。
6. *Executor*， *Executor* 是执行 *Action* 的组件，例如，一个 Python 的 *Action* 会被 一个 *Python Executor*  执行，对 *Executor* 的调用通过 HTTP RPC 进行。

## 部署

*breeze* 平台可以简便得在 *Docker Swarm* 上部署，使用 `deploy/` 目录下的脚本可以轻易得部署出一个 *breeze* 实例。但是需要考虑：

1. 通过 *Docker Swarm* 的 *Overlay Network* 通信有较大的性能损失，不管是在 *latency* 还是 *throughput* 上。
2. 目前 *Docker Images* 寄存在前公司的服务器上，随时可能不可用。

## 其他

感谢我的前公司（是谁我就不说啦！）授权我把 *breeze* 作为一个个人项目发布。如果这个项目对您有所帮助，那是我的荣幸。

另外，如果您希望对 *breeze* 进行商业使用，请通过 matao.xjtu#gmail.com 联系我。谢谢！

