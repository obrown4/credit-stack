import { signInWithPopup } from "firebase/auth";
import {  } from "./firebase";

export function SignIn(){
    signInWithGoogle = () => {
    const provider = new 
    auth.signInWithPopup(provider);
  }

  return (
    <>
      <button className="sign-in" onClick={signInWithGoogle}>Sign in with Google</button>
      <p>Do not violate the community guidelines or you will be banned for life!</p>
    </>
  )
}

export function SignOut(){
    return auth.currentUser && (
        <button className="sign-out" onClick={() => auth.signOut()}>Sign Out</button>
    )
}
