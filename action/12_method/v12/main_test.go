package v12

import "testing"

/*
TestEmployeeMaleCount 测试用例建立了一个 fakeStmtForMaleCount 的伪对象类型，然后在这个类型中嵌入了 Stmt 接口类型。
这样 fakeStmtForMaleCount 就实现了 Stmt 接口，我们也实现了快速建立伪对象的目的。
接下来我们只需要为 fakeStmtForMaleCount 实现 MaleCount 所需的 Exec 方法，就可以满足这个测试的要求了。
*/

type fakeStmtForMaleCount struct {
	Stmt
}

func (fakeStmtForMaleCount) Exec(stmt string, args ...string) (Result, error) {
	return Result{Count: 5}, nil
}

func TestEmployeeMaleCount(t *testing.T) {
	fakeStmt := fakeStmtForMaleCount{}
	count, _ := MaleCount(fakeStmt)
	if count != 5 {
		t.Errorf("want: %d, actual: %d", 5, count)
		return
	}
}
