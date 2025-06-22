# 单阶段流水线

```bash
pipeline {
    agent any 
    stages {
        stage('Stage 1') {
            steps {
                echo 'Hello world!' 
            }
        }
    }
}
```

# 多阶段流水线

```bash
pipeline {
    agent any 
    stages {
        stage('Stage 1') {
            steps {
                echo 'Hello world!' 
            }
        }        
        stage('Stage 2') {
            steps {
                echo 'This is Stage 2' 
            }
        }
    }
}
```

# 构建代理

```bash
agent any
agent none
agent { node { label 'labelName' } }
agent { label 'labelName' }
agent {
    node {
        label 'my-defined-label'
        customWorkspace '/some/other/path'
    }
}

pipeline {
    agent none
    stages {
        stage('Example Build') {
            agent { 
            	docker {
            		label 'slave'
            		'maven:3.9.3-eclipse-temurin-17' 
            		// 如果需要加参数
            		// args '-v /path/to/mount:/tmp'
            	}
            }
            steps {
                echo 'Hello, Maven'
                sh 'mvn --version'
            }
        }
        stage('Example Test') {
            agent { docker 'openjdk:17-jre' }
            steps {
                echo 'Hello, JDK'
                sh 'java -version'
            }
        }
    }
}
```

# options使用

```
// 全局配置
options {
    retry(3) 
	timeout(time: 1, unit: 'HOURS') //MINUTES, SECONDS
	timestamps()
	parallelsAlwaysFailFast()
	disableConcurrentBuilds(abortPrevious: true)
	buildDiscarder(logRotator(numToKeepStr: '10'))
	skipDefaultCheckout()
}

// 局部配置
pipeline {
    agent any 
    stages {
        stage('Stage 1') {
            steps {
                timeout(time: 3, unit: 'SECONDS'){
                    sh 'sleep 4'
                }
            }
        }        
        stage('Stage 2') {
            steps {
                retry(2){ sh 'echo hello' }
            }
        }
    }
}
```

# 环境变量

```
pipeline {
    agent any
    environment {
        CC = 'clang'  // 全局环境变量
    }
    stages {
        stage('Example') {
            // 局部环境变量，只能够在stage中使用
            environment {
               bb = 'golang'
            }			          
            steps {
            	withEnv(['name1=join','name2=jack']){ 
                	echo "my name is ${env.name1}"
                	sh 'my name is ${name2}'
                }
                script {
                      env.name3 = "peter"                 
                      env.name4 = "lisha"
                }
                echo "my name is ${env.name3}"
                sh 'my name is ${name4}'
            }
        }
    }
}
```

# 用户名和密码凭证

```bash
pipeline {
    agent any
    stages {
        stage('Example Username/Password') {
            environment {
                SERVICE_CREDS = credentials('my-predefined-username-password')
            }
            steps {
                sh 'echo "Service user is $SERVICE_CREDS_USR"'
                sh 'echo "Service password is $SERVICE_CREDS_PSW"'
                sh 'curl -u $SERVICE_CREDS https://myservice.example.com'
            }
        }
        stage('Example SSH Username with private key') {
            environment {
                SSH_CREDS = credentials('my-predefined-ssh-creds')
            }
            steps {
                sh 'echo "SSH private key is located at $SSH_CREDS"'
                sh 'echo "SSH user is $SSH_CREDS_USR"'
                sh 'echo "SSH passphrase is $SSH_CREDS_PSW"'
            }
        }
    }
}
```

# post介绍

```bash
pipeline {
    agent any
    stages {
        stage('Example') {
            steps {
                echo 'Hello World'
            }
        }
    }
    post { 
        always { 
            echo 'always'
        }
        success { 
            echo 'success'
        }
        failure { 
            echo 'failure'
        }
        aborted { 
            echo 'aborted'
        }
        unstable { 
            echo 'unstable'
        }        
    }
}
```

# 用户变量与清理命名空间

```bash
pipeline {
    agent any
    stages {
        stage('Example') {
            steps {
				script {		
              		cleanWs()  
          		}
            }
        }
    }
    post { 
        success {         	
            script {
                env.abc = currentBuild.durationString.split("and counting")[0]
                currentBuild.description = "构建用户 ${BUILD_USER} \n 持续时间 ${abc}"
                currentBuild.displayName = "部署项目: ${JOB_NAME} 部署URL: ${JOB_URL} 构建URL: ${BUILD_URL} 部署阶段: ${STAGE_NAME}"
            }
        }
    }
```

# 计划任务触发器

```bash
// Declarative //
pipeline {
    agent any
    triggers {
        cron('H */4 * * 1-5')
        //triggers{ cron('H H(9-16)/2 * * 1-5') }
    }
    stages {
        stage('Example') {
            steps {
                echo 'Hello World'
            }
        }
    }
}
```

# 并行执行

```bash
pipeline {
    agent any
    stages {
        stage('Non-Parallel Stage') {
            steps {
                echo 'This stage will be executed first.'
            }
        }
        stage('Parallel Stage') {
            parallel {
                stage('Branch A') {
                    steps {
                        echo "On Branch A"
                    }
                    }
                stage('Branch B') {
                    steps {
                        echo "On Branch B"
                    }
                }
            }
        }
    }
}	
```

# when

```bash
pipeline {
    agent any
    environment {
       DEPLOY_TO = 'production'
    }		    
    stages {
        stage('Example Build') {
            steps {
                echo 'Hello World'
            }
        }
        stage('Example Deploy') {
            when {
                anyOf {
                    environment name: 'DEPLOY_TO', value: 'production'
                    expression { BRANCH_NAME ==~ /(production|master)/ }
                    expression {
                        return params.BUILD.NUMBER >= 1
                    }
                }
            }
            steps {
                echo 'Deploying'
            }
        }
    }
}
pipeline {
    agent any
    parameters {
        choice choices: ['master', 'pre', 'test'], name: 'build_branch'
    }
    stages {
        stage('deploy') {
            when {
                not {
                    equals expected: build_branch, actual: 'master'
                }
            }
            steps {
                echo "deploy success"
            }
        }
    }
}
```

# input

```bash
pipeline {
    agent any
    stages {
        stage('Example') {
            input {
                message "这个构建你自已确认吗?"
                ok "是的，我确认."
                parameters {
                    string(name: 'PERSON', defaultValue: '我要发布', description: '你是为什么什么原因要构建呢?')
                }
            }
            steps {
                echo "Hello, ${PERSON}, nice to meet you."
            }
        }
    }
}
```

# Groovy postbuild

```bash
pipeline {
    agent any
    stages {
        stage('test') {
            steps {
                echo "success"
            }
        }
    }
    post {
        success {
            script {
                manager.addShortText("构建用户：${BUILD_USER}")
            }
        }
    }
}
```

# 更新流水线状态

```bash
pipeline {
    agent any
    options {
      gitLabConnection('your-gitlab-connection-name')
    }  
    stages {
      stage("build") {
        steps {
          updateGitlabCommitStatus name: 'build', state: 'running'
          echo "hello world"
        }
      }
    }    
    post {
      failure {
        updateGitlabCommitStatus name: 'build', state: 'failed'
      }
      success {
        updateGitlabCommitStatus name: 'build', state: 'success'
      }
      aborted {
        updateGitlabCommitStatus name: 'build', state: 'canceled'
      }
    }
}
```

# webhook触发

```bash
#!groovy

pipeline {
    agent any
    
    parameters {
      choice choices: ['refs/heads/pre', 'refs/heads/main', 'refs/heads/test'], name: 'branch_name'
      string defaultValue: 'http://10.0.7.30/jenkins/jenkinsfile.git', name: 'giturl'
    }
    
    triggers {
      GenericTrigger( 
          causeString: 'Generic Cause', 
          genericVariables: 
          [[defaultValue: '', key: 'branch_name', regexpFilter: '', value: '$.ref'],
          [defaultValue: '', key: 'giturl', regexpFilter: '', value: '$.project.git_http_url']], 
          printContributedVariables: true,
          printPostContent: true, 
          regexpFilterExpression: '',
          regexpFilterText: '', 
          token: 'abc123', 
          tokenCredentialId: ''
      )
    }

    stages {     
        stage('Stage 1') {
            steps {
                cleanWs()
                script {
                    branch = branch_name - 'refs/heads/'
                    //print "$branch"
                    checkout changelog: false, poll: false, scm: scmGit(
                    branches: [[name: "${branch}"]], extensions: [], 
                    userRemoteConfigs: [[credentialsId: 'gitlab', url: "${giturl}"]])
                    sh "cat README.md"
                }
            }
        }
    }
}
```

# 基于容器构建

```
pipeline {
    agent {
        node {
            label 'slave-docker'
        } 
    }
    environment {
		images_head = "registry.cn-hangzhou.aliyuncs.com"
    }  

    stages { 
        stage('Build') {
            agent {
                docker {
                    image 'maven:3.9.3-eclipse-temurin-17'
                    args '-v $HOME/.m2:/root/.m2'
                }
            }        
            steps {
                sh 'mvn;touch a.txt;touch /root/.m2/cache.txt'
            }
        }        
        stage('Hello') {
            agent {
                docker { 
                    label 'slave-docker'
                    image 'alpine:3.14' 
                }
            }
            steps {
                script {
                    sh """
                        ls -la
                        pwd
                        hostname
                        echo "#########################"
                    """
                }
            }
        }         
        stage('build-2') {          
            steps {
                script {
                    sh """
                    ls -la 
                    pwd
                    hostname
                    """
                    docker.withRegistry("https://${images_head}", 'aliyun-images-registry') {
                        def customImage = docker.build("${images_head}/tool-bucket/muke:${BUILD_TAG}-${GIT_COMMIT}")
                        customImage.push()                                              
                    }                        
                }                
            }
        }
                                        
    }
}
```

> ```
> // 	test-image从位于 的 Dockerfile构建./dockerfiles/test/Dockerfile
> def testImage = docker.build("test-image", "./dockerfiles/test")
> 
> // 可以docker build通过将其他参数添加到方法的第二个参数来传递其他参数build()。以这种方式传递参数时，字符串中的最后一个值必须是 docker 文件的路径，并且应以用作构建上下文的文件夹结尾
> def dockerfile = 'Dockerfile.test'
> def customImage = docker.build("my-image:${env.BUILD_ID}",
>                                    "-f ${dockerfile} ./dockerfiles")
> ```

# K8S动态节点

```
pipeline {
  agent {
    kubernetes {
     //inheritFrom 'mypod'
      yaml '''
        apiVersion: v1
        kind: Pod
        spec:
          containers:
          - name: maven
            image: maven:alpine
            command:
            - cat
            tty: true   
          - name: golang
            image: golang:1.16.5
            command:
            - sleep
            args:
            - 99d              
        '''
      retries 2
    }
  }
  stages {
    stage('') {
      steps {
        container('maven') {
          sh 'touch a.txt'
        }
      }
    }
    stage('Run golang') {
      steps {
        container('golang') {
          sh 'ls;pwd'
        }
      }
    }    
  }
}
```

# sonar接入ldap

```bash
# LDAP configuration
# General Configuration
sonar.security.realm=LDAP
ldap.url=ldap://10.0.7.30:389
ldap.bindDn=cn=admin,dc=muke,dc=cn
ldap.bindPassword=muke6666
  
# User Configuration
ldap.user.baseDn=ou=user,dc=muke,dc=cn
ldap.user.request=(&(objectClass=inetOrgPerson)(uid={login}))
ldap.user.realNameAttribute=cn
ldap.user.emailAttribute=mail
```
