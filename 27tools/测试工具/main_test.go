package main
/*
suite包提供了类似rails minitest中可以给每个测试用例进行前置操作和后置操作的功能，这个方便的功能，在前置操作和后置操作中去初始化和清空数据库，就可以帮助我们实现第一个目标。
同时，还可以声明在这个测试用例周期内都有效的全局变量
*/ 
import (
  "testing"
  "github.com/stretchr/testify/assert"
  // "github.com/stretchr/testify/require"
 )

// func TestCase1(t *testing.T) {
//     name := "Bob"
//     age := 10

//     assert.Equal(t, "bob", name)
//     assert.Equal(t, 20, age)
// }
func TestCase2(t *testing.T) {
  // assert.Equal(t, 123, 123, "they should be equal")
  // assert.NotEqual(t, 123, 456, "they should not be equal")
  //断言为nil(适用于错误)
  assert.Nil(t, nil)
  // 断言为not nil(当您期望某些东西时很好)
 assert.NotNil(t, "")
  // name := "Bob"
  // age := 10
  // require.Equal(t, "bob", name)
  // require.Equal(t, 20, age)
}