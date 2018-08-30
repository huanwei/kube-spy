# Netem参数说明
> 本文主要内容来自Linux基金会Wiki网站Netem文档，点击[这里](https://wiki.linuxfoundation.org/networking/netem)访问原文



netem通过模拟广域网的特性为测试协议提供[网络仿真](http://en.wikipedia.com/wiki/Network_emulation)功能。当前版本模拟可变延迟，丢失，重复和重新排序。

如果您运行当前的2.6发行版（[Fedora](http://en.wikipedia.com/wiki/Fedora_Core)，[OpenSuse](http://en.wikipedia.com/wiki/Open_Suse)，[Gentoo](http://en.wikipedia.com/wiki/Gentoo_Linux)，[Debian](http://en.wikipedia.com/wiki/Debian)，[Mandriva](http://en.wikipedia.com/wiki/Mandriva)，[Ubuntu](http://en.wikipedia.com/wiki/Ubuntu_%28Linux_distribution%29)），那么netem已在内核中启用，并且包含当前版本的[iproute2](https://wiki.linuxfoundation.org/networking/iproute2)。 netem内核组件在以下位置启用：

```
 Networking -->
   Networking Options -->
     QoS and/or fair queuing -->
        Network emulator
```
Netem由命令行工具`tc`控制，它是iproute2工具包的一部分。 tc命令使用`/usr/lib/tc`目录中的共享库和数据文件。
## 目录
* 模拟广域网延迟
* 指定延迟的分布
* 数据包丢失
* 数据包重复
* 数据包损坏
* 数据包重新排序
* HZ的值如何影响Netem？

## 模拟广域网延迟
这是最简单的示例，它只是为从本地以太网发出的所有数据包添加了固定数量的延迟。

`delay,100ms`

现在，在本地网络上进行主机的简单ping测试应显示增加100毫秒。延迟受内核（HZ）的时钟分辨率限制。在大多数内核版本为2.4的系统中，系统时钟以100hz运行，允许延迟增量为10ms。在2.6内核上，该值是1000到100赫兹的可配置参数。

真正的广域网具有可变性，因此可以添加随机变化。

`delay,100ms,10ms`

这导致增加的延迟为100ms±10ms。网络延迟变化不是纯粹随机的，因此要模拟存在[相关性](http://en.wikipedia.com/wiki/correlation)。

`delay,100ms,10ms,25%`

这导致增加的延迟为100ms±10ms，下一个随机元素取决于最后一个随机元素的25％。这不是真正的统计学意义上的相关性，而只是程序模拟的近似值。

## 指定延迟的分布
通常，网络中的延迟是不均匀的。使用类似[正态分布](http://en.wikipedia.com/wiki/Normal_Distribution)的东西来描述延迟的变化更为常见。 netem可以采用预先编写的表格来指定非均匀的分布。

`delay,100ms,20ms,distribution,normal`

实际使用的表格（normal，pareto，paretonormal）作为[iproute2](https://wiki.linuxfoundation.org/networking/iproute2)编译的一部分生成并放在`/usr/lib/tc`中；因此用户可以根据实验数据花费一点时间编写自己的分布表格。

## 数据包丢失
随机的数据包丢失在`tc`命令中以百分比形式指定。最小的非零值是：

1/（2的32次方） = 0.0000000232％

`loss,0.1%`

这导致0.1个百分点的（即1000个中的1个）数据包被随机丢弃。

还可以添加可选的相关性。使随机数不那么随机，可以用于模拟分组突发丢失。

`loss,0.3%,25%`

这将导致0.3％的数据包丢失，并且下一个包丢失的概率有四分之一取决于上一个数据包的丢包概率。

`丢包概率(n) = 0.25 * 丢包概率(n-1) + 0.75 * 指定概率`

## 数据包重复
数据包重复的指定方式与数据包丢失的方式相同。

`duplicate,1%`

## 数据包损坏
可以制造数据包损坏来模拟随机噪声（在2.6.16内核或更高版本中）。这个功能会在数据包中的随机位置引入单个比特的错误（翻转）。

`corrupt,0.1%`

## 数据包重新排序
指定数据包重新排序有两种不同的方法。第一个方法是`gap`，`gap`使用固定序列并重新排序每第N个数据包。一个简单的用例是：

`gap,5,delay,10ms`

这会导致每隔5（10,15，...）个数据包会立即发送，其他每个数据包都会延迟10ms。这种方式模拟的重排是可预测的。

重新排序的第二种方式更贴近真实网络。它会导致一定比例的数据包被错误排序：

`delay,10ms,reorder,25%,50%`

在此示例中，25％的数据包（相关性为50％）将立即发送，其他数据包将延迟10毫秒。

如果相邻数据包的随机延迟值不同，netem的较新版本也将重新排序数据包。以下参数将导致这一类的重新排序：

`delay,100ms,75ms`

如果第一个数据包的随机延迟为100毫秒（100毫秒基数 -  0毫秒抖动），第二个数据包在1毫秒后发送，延迟时间为50毫秒（100毫秒基数 -  50毫秒抖动），第二个数据包将首先发送。这是因为在netem内部的队列规则tfifo会按时间顺序进行数据包发送。

### 注意
* 混合形式的重新排序可能会导致意外结果；
* 任何方式的重新排序，都会引入延迟；
* 如果延迟小于数据包之间到达时间的时间差，则无法观察到数据包的重新排序。

## HZ的值如何影响Netem？

在2.6行内核中，HZ是一个可配置参数，取值为100,250或1000。因为它会影响Netem能够延迟数据包的粒度，所以最好将HZ设置为1000，允许以1ms的增量延迟。有关HZ影响的更详细讨论，请参阅[此邮件列表](http://lists.osdl.org/pipermail/netem/2006-March/000343.html)上的帖子。

在内核版本2.6.22或更高版本中，netem将使用高分辨率计时器（如果已启用）。这允许更精细的粒度（亚单位）分辨率。

