pipeline {
    agent any
    
    environment {
        IMAGE_NAME = 'registry.gitlab.com/hyperneutr0n/chirpy' 
        
        REPO_URL = 'https://github.com/hyperneutr0n/chirpy.git'
        
        REGISTRY_URL = 'registry.gitlab.com'
        DOCKER_CREDS = 'gitlab-container-registry-auth'
        GITHUB_CREDS = 'github-private-auth'
        KUBECONFIG_ID = 'k3s-kubeconfig'
        
        NAMESPACE = 'chirpy'
        DEPLOYMENT_NAME = 'chirpy-api'
    }

    stages {
        stage('Checkout') {
            steps {
                git credentialsId: "${GITHUB_CREDS}", url: "${REPO_URL}", branch: 'main'
            }
        }

        stage('Test') {
            steps {
                sh 'go test ./...'
            }
        }

        stage('Debug Info') {
            steps {
                echo "TAG_NAME is: ${env.TAG_NAME}"
                echo "GIT_BRANCH is: ${env.GIT_BRANCH}"
            }
        }

        stage('Build & Push Image') {
            when { 
                expression { return env.GIT_BRANCH?.contains('tags/') }
             }
            steps {
                script {
                    def ACTUAL_TAG = env.GIT_BRANCH.split('/').last()

                    docker.withRegistry("https://${REGISTRY_URL}", "${DOCKER_CREDS}") {
                        def customImage = docker.build("${IMAGE_NAME}:${TAG_NAME}")
                        
                        customImage.push()
                        
                        customImage.push("latest")
                    }
                }
            }
        }

        stage('Deploy to K3s') {
            when { 
                expression { return env.GIT_BRANCH?.contains('tags/') }
             }
            steps {
                def ACTUAL_TAG = env.GIT_BRANCH.split('/').last()
                
                withKubeConfig([credentialsId: "${KUBECONFIG_ID}"]) {
                    
                    sh "kubectl set image deployment/${DEPLOYMENT_NAME} ${DEPLOYMENT_NAME}=${IMAGE_NAME}:${TAG_NAME} -n ${NAMESPACE}"
                    
                    sh "kubectl rollout status deployment/${DEPLOYMENT_NAME} -n ${NAMESPACE}"
                }
            }
        }
    }
}