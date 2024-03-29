+++
title = "常见算法思想"
date = "2021-05-27T14:32:31+08:00"
author = "vincixu"
authorTwitter = "" #do not include @
cover = ""
tags = ["app-design"]
keywords = ["algorithm"]
description = ""
showFullContent = false
+++

# 常见算法思想

算法几个重要思想：
- `枚举`: 最简单的暴力方法就是枚举所有可能项
- `回溯`：思想与枚举类似，尝试所有的可能项，不同点在于它在达到临界条件后会退回到某个分叉点，是一个递归的思想，整个回溯的过程可以看作对树的遍历(和前序中序无关，那是二叉树的遍历方法)，所以可以做剪枝优化
- `分治`：将一个问题可以拆解成规模更小的子问题，从而各个击破。
- `贪心`：每次选择都选择最优项，是一种不考虑后继步骤的思想
- `动态规划`：将问题拆解规模更小的子问题来解决，和分治类似，最大不同点在于，动态规划需要满足
  + 最优子结构：最优子结构的意思是局部最优解能决定全局最优解（对有些问题这个要求并不能完全满足，故有时需要引入一定的近似）。简单地说，问题能够分解成子问题来解决。
  + 无后效性：即子问题的解一旦确定，就不再改变，不受在这之后、包含它的更大的问题的求解决策影响。
  + 重叠子问题：子问题重叠性质是指在用递归算法自顶向下对问题进行求解时，每次产生的子问题并不总是新问题，有些子问题会被重复计算多次。动态规划算法正是利用了这种子问题的重叠性质，对每一个子问题只计算一次，然后将其计算结果保存在一个表格中，当再次需要计算已经计算过的子问题时，只是在表格中简单地查看一下结果，从而获得较高的效率，降低了时间复杂度
- `滑动窗口`：滑动窗口可用于解决最优 `连续子集/子序列` 的问题，比如最大不重复子串
- `位运算`：可用于解决比如组合，存储结果集等
- `快慢指针`：两个速度不同的指针，用于检测是否存在环


几个重要的数据结构：
- `栈（Stack）`:  先进后出，可以用来处理括号匹配，逆序等。可以用链表 or 数组实现。
- `队列（Queue）`：先进先出，可以用来处理广度搜索。可以用链表 or 数组实现。队列又有优先队列一说，即某高优项可以显出，通过建堆的方法可以在 O(logn) 的复杂度下实现
- `数组（Array）`：跳过, 值得一提的是桶和这个东西类似
- `树（Tree）`
  + `多叉树`：b, b-, b+ 树
  + `二叉树`：二叉排序树(BinarySearchTree)，其中又可以分为常见 红黑树 和 AVL 树，AVL维护插入复杂度比红黑树高，但是搜索效率比红黑高。
- `堆（Heap）`：分为大顶堆和小顶堆，可以看作特殊的完全二叉树，建堆可以通过：1、插入新元素，2、交换元素(与父元素相比，满足条件则交换)，复杂度 O(nlogn)。bottom-up 建堆方法是对所有非叶子节点（n/2 - 1）进行局部堆化，由底至顶，算法复杂度在线性时间。
- `散列表（Hash table）`: 又称 Hash表，指通过哈希函数将 key 映射到数组的序号中，关键问题在于解决冲突，解决冲突的方法：开放寻址（线性搜索、再哈希）、拉链法（使用链表 or 桶存储）
- `前缀/后缀 集`:
- `并查集`：并查集是为了解决关联性而产生的数据结构，比如需要对数组进行分类