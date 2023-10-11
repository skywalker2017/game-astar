# 对战接口说明

## 枚举类
### 单位移动类型枚举 MoveType
+ stand 不移动 0
+ moveGround 地面移动 1
+ moveFly 飞行 2
+ moveBounce 跳跃 3
+ moveDig 地下 4

### 单位生存状态枚举 LivingStatus
+ destroyed 摧毁 -1
+ unDeployed 未部署 0
+ deployed 已部署 1

### 单位攻击状态枚举 AtkStatus
+ free 自由 0
+ searching 搜索目标中 1
+ moving 移动中 2
+ attacking 攻击中 3
+ wallAttacking 攻击围墙 4 （攻击围墙状态的单位可随时触发搜索目标）

### 防御单位类型枚举 DefenderType
+ wall 围墙 1
+ trainStructure 训练建筑 2
+ resource 资源建筑 3
+ groundDefense 地面防御 4
+ airDefense 空中防御 5

### 攻击类型枚举 AtkType
+ normalAtk 普通攻击 0
+ rangeAtk 溅射攻击 1

## battle方法

### InitBattle(size, buildingSize int) int
+ 功能： 初始化对战
+ 参数：
  + size: 地图大小
  + buildingSize：防御建筑数目
+ 返回值：
  + battle index
### GetBattle(index int) *Battle
+ 功能： 获取对战对象
+ 参数：
  + index: battle index
  + buildingSize：防御建筑数目
+ 返回值：
  + battle 对象
### AddAttacker(x, y, atkRange, damage, deathDamage int, atkTypeInt int, atkPriorityInt int, suicideInt int) int
+ 功能： 加载攻击单位
+ 参数：
  + x: x坐标
  + y: y坐标
  + atkRange: 攻击距离
  + damage: 单位时间伤害
  + deathDamage: 死亡时伤害
  + atkTypeInt: AtkType
  + atkPriorityInt: DefenderType
  + suicideInt: 0-普通 1-自杀攻击 
+ 返回值：
  + 对应battle中的 attacker index 
### AddDefender(xMin, yMin, dSize, damage, atkRange, health int, defenderType int) int
+ 功能： 加载防御单位
+ 参数：
  + xMin: 左上角x坐标
  + yMin: 左上角y坐标
  + dSize: 实体大小
  + damage: 单位时间伤害
  + atkRange: 攻击距离
  + health: 生命值
  + defenderType: DefenderType
+ 返回值：
  + 对应battle中的 defender index

### GetAttacker(attackerIndex int) Attacker
+ 功能： 获取对战攻击单位
+ 参数：
  + attackerIndex: attacker index
+ 返回值：
  + 对应battle中attacker对象

### GetDefender(defenderIndex int) Defender
+ 功能： 获取对战防御单位
+ 参数：
  + defenderIndex: defender index
+ 返回值：
  + 对应battle中defender对象

### Play()
+ 功能： 对战状态更新
+ 参数：
+ 返回值：

## Attacker方法

### Move()
+ 功能：执行移动操作

### Attack()
+ 功能：执行攻击操作

### Search()
+ 功能：执行巡路操作

### GetIndex() int
+ 功能：获取attacker index
+ 返回值：
  + 对应battle中attacker index

### GetBattle() Battle {
+ 功能：获取对战
+ 返回值：
  + battle对象

### GetTarget() int {
+ 功能：获取当前目标
+ 返回值：
  + 对应battle中defender index

### GetLivingStatus() int {
+ 功能：获取生存状态
+ 返回值：
  + LivingStatus

### GetAtkStatus() int {
+ 功能：获取攻击状态
+ 返回值：
  + AtkStatus

### GetSubPos() *Point
+ 功能：获取子坐标
+ 返回值：
  + x, y

### GetPos() *Point
+ 功能：获取坐标
+ 返回值：
  + x, y


## Defender方法

### Attack()
+ 功能：执行攻击操作

### Search()
+ 功能：执行巡路操作

### GetIndex() int
+ 功能：获取defender index
+ 返回值：
  + 对应battle中defender index

### GetBattle() Battle {
+ 功能：获取对战
+ 返回值：
  + battle对象

### GetTarget() int {
+ 功能：获取当前目标
+ 返回值：
  + 对应battle中attacker index

### GetLivingStatus() int {
+ 功能：获取生存状态
+ 返回值：
  + LivingStatus

### GetAtkStatus() int {
+ 功能：获取攻击状态
+ 返回值：
  + AtkStatus

### GetXMin() int {
+ 功能：获取最小x坐标
+ 返回值：
  + 坐标

### GetXMax() int {
+ 功能：获取最大x坐标
+ 返回值：
  + 坐标

### GetYMin() int {
+ 功能：获取最小y坐标
+ 返回值：
  + 坐标

### GetYMax() int {
+ 功能：获取最大y坐标
+ 返回值：
  + 坐标

