#define _GNU_SOURCE
#include <sys/mount.h>
#include <sys/types.h>
#include <sys/wait.h>
#include <stdio.h>
#include <sched.h>
#include <signal.h>
#include <unistd.h>
#define STACK_SIZE (1024 * 1024)
static char container_stack[STACK_SIZE];
char* const container_args[] = {
  "/bin/bash",
  NULL
};

int container_main(void* arg)
{
  printf("Container - inside the container!\n");
  execv(container_args[0], container_args);
  // 如果物理机挂载类型是shared，那么需要重新挂载根目录
  // mount("","/",NULL,"MS_PRIVATE","")
  // 通过重新挂载/tmp目录去覆盖Mount Namespace中挂载内容与宿主机一样的问题
  // tmpfs 内存盘格式
  mount("none","/tmp","tmpfs",0,"")
  printf("Something's wrong!\n");
  return 1;
}

int main()
{
  printf("Parent - start a container!\n");
  int container_pid = clone(container_main, container_stack+STACK_SIZE, CLONE_NEWNS | SIGCHLD , NULL);
  waitpid(container_pid, NULL, 0);
  printf("Parent - container stopped!\n");
  return 0;
}