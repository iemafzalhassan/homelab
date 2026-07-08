pipeline {
    agent {
        label 'kaniko'
    }
    stages {
        stage('Build with Kaniko') {
            steps {
                container('kaniko') {
                    sh '''
                        echo "FROM alpine:latest" > Dockerfile
                        echo "CMD echo 'Hello from Jenkins Spot Kaniko Build!'" >> Dockerfile
                        
                        # We would normally push to ACR, but first we need ACR credentials
                        # For Phase 7 validation, we just ensure Kaniko can execute
                        /kaniko/executor --context `pwd` --dockerfile `pwd`/Dockerfile --no-push
                    '''
                }
            }
        }
    }
}
