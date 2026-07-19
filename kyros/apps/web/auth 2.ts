import NextAuth from "next-auth"
import Keycloak from "next-auth/providers/keycloak"

export const { handlers, signIn, signOut, auth } = NextAuth({
  providers: [
    Keycloak({
      clientId: process.env.KEYCLOAK_CLIENT_ID || "kyros-web",
      clientSecret: process.env.KEYCLOAK_CLIENT_SECRET || "kyros-web-secret-for-dev",
      issuer: process.env.KEYCLOAK_ISSUER || "http://localhost:8081/realms/kyros",
    })
  ],
  callbacks: {
    async jwt({ token, account, user }) {
      if (account && user) {
        token.accessToken = account.access_token;
        token.idToken = account.id_token;
        
        // Sync user to Go API on first login/token refresh
        try {
          const apiUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
          await fetch(apiUrl + "/v1/auth/sync", {
            method: "POST",
            headers: {
              "Authorization": `Bearer ${account.access_token}`,
              "Content-Type": "application/json"
            }
          });
        } catch (e) {
          console.error("Failed to sync user to Go API", e);
        }
      }
      return token;
    },
    async session({ session, token }) {
      // @ts-ignore
      session.accessToken = token.accessToken as string;
      return session;
    }
  },
  pages: {
    signIn: "/login",
  }
})
