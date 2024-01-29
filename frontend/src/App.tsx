import { useEffect, useState } from "react";
import { Auth } from "@supabase/auth-ui-react";
import { User, createClient } from "@supabase/supabase-js";
import { ThemeSupa } from "@supabase/auth-ui-shared";

const supabaseProjectId = import.meta.env.VITE_SUPABASE_PROJECT_ID;
const supabaseUrl = `https://${supabaseProjectId}.supabase.co`;
const publicAnonKey = import.meta.env.VITE_SUPABASE_ANON_KEY;

const supabase = createClient(supabaseUrl, publicAnonKey);

function App() {
  // Contains the user received from Supabase
  const [user, setUser] = useState<User>();

  useEffect(() => {
    supabase.auth.onAuthStateChange((event, session) => {
      if (event === "SIGNED_IN") {
        setUser(session?.user);
      }
    });
  }, []);

  return user ? <LoggedIn /> : <LoggedOut />;
}

function LoggedIn() {
  const [ourSecretData, setOutSecretData] = useState();

  useEffect(() => {
    fetch("http://localhost:3000/secret", {
      method: "POST",
      headers: {
        Authorization: getToken(),
      },
    })
      .then((res) => res.json())
      .then((data) => setOutSecretData(data));
  }, []);

  const handleSignOut = () => {
    supabase.auth.signOut().then(() => {
      window.location.reload();
    });
  };

  return (
    <>
      <div>{JSON.stringify(ourSecretData)}</div>
      <button onClick={handleSignOut}>Sign out</button>
    </>
  );
}

function LoggedOut() {
  return (
    <Auth
      supabaseClient={supabase}
      appearance={{
        theme: ThemeSupa,
      }}
      providers={[]}
      theme="dark"
      redirectTo="/"
      showLinks
    />
  );
}

const getToken = () => {
  const storageKey = `sb-${supabaseProjectId}-auth-token`;
  const sessionDataString = localStorage.getItem(storageKey);
  const sessionData = JSON.parse(sessionDataString || "null");
  const token = sessionData?.access_token;

  return token;
};

export default App;
