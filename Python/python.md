# Python

### 保留关键字

```python
import keyword

print(keyword.kwlist)

# pass 占位符
```

### 内置函数

```python
number = "1"

# 内存地址
print(id(number))
print(type(number))

# range 产生一个生成器，左闭右开，起始值和步进可以省略
print(list(range(0, 10, 2)))

num = input("please input a num: ")
print(num)
# 默认str类型
print(type(num))

# 3, 元素在列表中的首个index，后面参数是查询的起始和结束的位置，结束位置可超出索引，均可省略，没有结果抛异常
print("abcaef".index("ae", 2, 9))

# 可以同时删除多个，且是不同类型
del ...

# 排序
a = ["a", "b", "c"]
sorted(a, reverse=True)
print(a)
# ['a', 'b', 'c'] 内置函数会产生新的list

a.sort(reverse=True)
print(a)
# list的方法是直接修改自己，不产生新的list
['c', 'b', 'a']

# dir 查看对象有哪些方法，可以查看类私有方法，然后通过_Class__methond进行调用
dir()
```

### 变量

```python
name = "Jack"
print(id(name))

name= "Jinke"
print(id(name))
# 会指向新的内存地址，产生内存垃圾
```

### 浮点数

```python
import decimal

# 十进制转二进制底层问题
print(1.1 + 2.1)
# 3.2
print(1.1 + 2.2)
# 3.3000000000000003 
print(decimal.Decimal("1.1") + decimal.Decimal("2.2"))
# 3.3
```

### 类型转换

```python
number = "1.1"

# 三者可以相互转换
print(type(int(number)))
print(type(float(number)))
print(type(str(int(number))))

a: int = 1
b: float = 2.3
c: str = "3"
# float
print(type(a + b))
# error
print(type(a + c))
```

### 运算符

```python
print(11 / 2) # 除法
print(11 // 2) # 向下取整 5
print(-11 // 2) # -6
print(11 % 2) # 取余
print(2**2) # 幂运算

# 赋不同类型的值
a, b, c = 1, "a", True
print(a, b, c)

# 交换变量值，无须中间变量
a, b = b, a
print(a, b)

# 比较值和内存地址
a, b = 1, 1
print(a == b) # 值相等
print(a is b) # 内存地址一样

# in 运算符
string = "abc"
print("a" in string)
print("a" not in string)

fruit = ["apple", "orange", "banana"]
print("apple" in fruit)
```

### 内存模型

```python
a = b = (2, 3, 4)
c = (2, 3, 4)
# a,b内存地址一样,c 不一样

# pycharm 做了优化，跟shell不一样
# [-5, 256] ,只存储包含标准字符（数字、字母、下划线）的字符串会驻留内存
# 其他普通类型对象的新建，python会请求内存，申请内存 。当n1的引用指向其他对象时，原有对象的引用计数会自动减1，没有被引用的对象会立即回收。

# 内存回收三种方式
# 1.引用计数
# 2.标记清除，打破了循环引用
# 3.分代回收，活的越长的对象，就越不可能是垃圾
```

### 循环

```python
# while...else
# for...else
# 循环正常结束，中间没有break会执行else语句
for i in range(3):
    print(i)
else:
    print("else")
```

### 切片

```python
a = ["a", "b", "c"]
# [start:end:step] 左闭右开，拷贝会申请新的内存，step为-1时会revert list
b = a[:]
print(id(a), a, id(b), b)
```

### 列表

```python
a = ["a", "b", "c"]
b = ["x", "y", "z"]

a[2:] = b
print(a)
['a', 'b', 'x', 'y', 'z']

# 列表生成式
a = list(i for i in range(1, 5))
# a = [i for i in range(1, 5)]
print(a)
```

### 字典

```python
# 创建字典的两种方式
a = {"a": "A"}
b = dict(a="a")
print(a, b)
# 没有该值，结果为None
print(a.get("b"))
# get可以设置默认值
print(a.get("b", "B"))
# 判断是否有该值
print("a" in a)
# 这个直接抛异常
print(a["b"])

# 字典的视图，keys,values,items
a = {"a": "A"}
print(type(a.keys()), type(a.values()), type(a.items()))
# <class 'dict_keys'> <class 'dict_values'> <class 'dict_items'>
for item in a.items():
    # 元组类型
    print(type(item))
    print(item)
    
# 字典生成式，以元素少的list为元素个数
a = ["a"]
A = ["A"]
z = zip(a, A)
d = dict(z)
g = {key: value for key, value in z}
print(type(z), z)
print(type(d), d)
print(type(g), g)
# <class 'zip'> <zip object at 0x103505900>
# <class 'dict'> {'a': 'A'}
# zip的后值只能取一次，之后会被清空
# <class 'dict'> {}
```

### 集合

```python
# 是没有key的字典，无序且不重复，可以和元组、列表互转
a = {"a", "b"}
b = set(("a", "b", "a"))
print(type(a), a)
print(type(b), b)

# 一次增加一个元素
a.add("c")
# 一次增加多个元素，iterable类型
a.update(range(5))

# 没有该元素会抛异常
a.remove("d")
# 不会抛异常
a.discard("d")

# 集合生成式
a = set(i for i in range(5))
# a = {i for i in range(5)}
print(a)
```

### 元组

```python
a = ("a", "b")
# 多个元素可以省略括号
b = "a", "b"
# 只有一个元素时需要加","，否则时string类型
c = ("a",)
# 里面有括号的嵌套
d = tuple(("a", "b"))
print(a, b, c, d)

# 元组的元素不可修改，元素如果是可变类型，元素的元素可以修改
```

### 函数

```python
# 参数是不可变类型会拷贝值到函数，可变类型会传地址值，可以读取全局变量，想要修改全局变量需要声明该变量为global
a = 1
b = [2]

def test(x: int, y: [int]):
    x = 3
    a = 2
    y.append(x)

test(a, b)
print(a, b)
# 1 [2, 3]

# 返回值大于1个时，会返回元组
def test():
    return 1, 2

print(test())

# 默认参数
def test(a=10):
# 可变位置参数，元组类型
def test(*args):
# 可变关键字参数，字典类型
def test(**args):
```

### 异常

```python
try:
    pass
except Exception as e:
    pass
else:
    pass
finally:
    pass
```

### 面向对象

```python
class Student:
    name = "jack"  # 类属性，所有对象都可访问，更改后所有的对象都能看到改后的值

    def __init__(self, age):
        self.age = age  # 将局部变量赋值给实例属性

    # 实例方法
    def a(self):
        pass

    # 静态方法
    @staticmethod
    def b():
        pass

    # 类方法
    @classmethod
    def c(cls):
        pass

# stu 是实例对象，Student是类对象   
stu = Student(age=28)

# 动态绑定
# 类可以动态绑定属性，所有创建的实例都可访问
# 实例可动态绑定属性和方法，不会影响类和其它实例

# 继承，python支持多继承，但坑太多，不要用
class A:
    def __init__(self, a):
        self.a = a

class B(A):
    def __init__(self, a, b):
      	# init方法得先调用父类的init
        super().__init__(a)
        self.b = b
        
# 子类可以重写父类方法，只要方法名一样就是重写

# 因为所有类都继承了object，所以可以重写object提供的私有方法，比如__str__,__gt__,__hash__

# __new__ 是创建对象，这个阶段已经通过调用object的__new__方法分配好了内存，并返回新创建的对象到__init__(self)中，__init__ 是初始化对象
```

