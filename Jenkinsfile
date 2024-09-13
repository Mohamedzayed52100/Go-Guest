properties([pipelineTriggers([githubPush()])])

pipeline {
    agent {
        kubernetes {
            yaml '''
apiVersion: v1
kind: Pod
spec:  
  serviceAccountName: jenkins-agent
  containers:
  - name: shell
    image: 348277991470.dkr.ecr.eu-west-2.amazonaws.com/jenkins-agent:latest
    command:
    - sleep
    args:
    - infinity
  - name: docker-dind
    image: 348277991470.dkr.ecr.eu-west-2.amazonaws.com/jenkins-dind:latest
    securityContext:
      privileged: true
    volumeMounts:
    - name: dind-storage
      mountPath: /var/lib/docker
  tolerations:
  - key: "ENV"
    operator: "Equal"
    value: "JENKINS"
    effect: "NoSchedule"
  nodeSelector:
    role: eks-jenkins-node-groups        
  volumes:
  - name: dind-storage
    emptyDir: {}
'''
            defaultContainer 'shell'
        }
    }

    options {
        disableConcurrentBuilds()
    }

    environment {
        IMAGE_REGISTRY_PROD="348277991470.dkr.ecr.eu-west-2.amazonaws.com/goplace-guest-prod"
        IMAGE_REGISTRY_DEV="348277991470.dkr.ecr.eu-west-2.amazonaws.com/goplace-guest-dev"
        IMAGE_REGISTRY_STAGING="348277991470.dkr.ecr.eu-west-2.amazonaws.com/goplace-guest-staging"
        K8S_MANIFESTS_BRANCH_DEV="goplace-guest-dev"
        K8S_MANIFESTS_BRANCH_PROD="goplace-guest-prod"
        K8S_MANIFESTS_BRANCH_STAGING="goplace-guest-staging"
    }

    stages {

        stage('Setting-up variables') {
            steps {
                script {
                    env.TAG = sh(
                        returnStdout: true,
                        script: '''#!/bin/bash
                            git config --global --add safe.directory "*"
                            COMMIT_COUNT=$(git rev-list --count HEAD)
                            SHORT_HASH=$(git rev-parse --short HEAD)
                            echo ${COMMIT_COUNT}.${SHORT_HASH}
                        '''
                  ).trim()
                    echo "The tag is ${env.TAG}"
                    if (env.GIT_BRANCH == 'dev') {
                        env.IMAGE_REGISTRY = env.IMAGE_REGISTRY_DEV
                        env.K8S_MANIFESTS_BRANCH = env.K8S_MANIFESTS_BRANCH_DEV
                        env.DOCKERFILE_NAME = 'Dockerfile.dev'
                    } else if (env.GIT_BRANCH == 'main'){
                        env.IMAGE_REGISTRY = env.IMAGE_REGISTRY_PROD
                        env.K8S_MANIFESTS_BRANCH = env.K8S_MANIFESTS_BRANCH_PROD
                        env.DOCKERFILE_NAME = 'Dockerfile'
                    } else if (env.GIT_BRANCH == 'staging'){
                        env.IMAGE_REGISTRY = env.IMAGE_REGISTRY_STAGING
                        env.K8S_MANIFESTS_BRANCH = env.K8S_MANIFESTS_BRANCH_STAGING
                        env.DOCKERFILE_NAME = 'Dockerfile'
                    } else {
                        currentBuild.result = 'FAILURE'
                        error("Branch is not configured")
                    }
                }
            }
        }

        stage('Pulling K8S Deployment Files') {
            steps {
                dir('k8s') {
                    git (url: "https://github.com/goplaceapp/k8s-manifests.git", credentialsId: "Github-Credintials", branch: "${env.K8S_MANIFESTS_BRANCH}")
                }
            }
        }
        
stage('Build and Push to AWS ECR') {
    steps {
        container("docker-dind") {
            script {
                withCredentials([string(credentialsId: 'GitHubAccessToken', variable: 'GITHUB_ACCESS_TOKEN')]) {
                    sh """ 
                        aws ecr get-login-password --region eu-west-2 | docker login --username AWS --password-stdin 348277991470.dkr.ecr.eu-west-2.amazonaws.com
                        docker build -t \${IMAGE_REGISTRY}:\${TAG} -f \${DOCKERFILE_NAME} --build-arg GITHUB_ACCESS_TOKEN=\${GITHUB_ACCESS_TOKEN} .
                        docker push \${IMAGE_REGISTRY}:\${TAG}
                    """
                }
            }
        }
    }
}

        stage('Deploy') {
            steps {
                sh '''
                export COMMIT_COUNT=$(git rev-list --count HEAD)
                export SHORT_HASH=$(git rev-parse --short HEAD)
                export TAG=${COMMIT_COUNT}.${SHORT_HASH}
                sed -i.bak "s|IMAGE_NAME|${IMAGE_REGISTRY}:${TAG}|g" k8s/deployment.yaml
                kubectl delete -f k8s/secret.yaml --ignore-not-found
                kubectl apply -f k8s/secret.yaml
                kubectl apply -f k8s/deployment.yaml
                kubectl apply -f k8s/service.yaml
                '''
            }
        }
    }
}
