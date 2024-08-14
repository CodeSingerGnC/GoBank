import subprocess

def run_make_test(times):
    for i in range(1, times + 1):
        print(f"执行次数: {i}")
        result = subprocess.run(["make", "test"], capture_output=True, text=True)
        
        if result.returncode != 0:
            print(f"测试在第 {i} 次执行时失败。")
            print("错误输出:")
            print(result.stderr)
            return False
    print("所有测试均成功完成。")
    return True

if __name__ == "__main__":
    if not run_make_test(100):
        exit(1)