import jenkins.model.*
import hudson.security.*

def jenkins = Jenkins.getInstance()
def hudsonRealm = jenkins.getSecurityRealm()

if (hudsonRealm instanceof HudsonPrivateSecurityRealm) {
    hudsonRealm.createAccount('admin', 'admin123')
    jenkins.save()
    println "Admin password has been reset to admin123"
}
