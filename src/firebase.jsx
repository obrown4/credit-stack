// firebase sdk
import { initializeApp } from "firebase/app";
import { getAuth } from "firebase/auth";
import { getFirestore } from "firebase/firestore"



const app = initializeApp({
  apiKey: "AIzaSyAyTrR6thdjAQRXbOGM5C9SAzpEl9q7TRY",
  authDomain: "credit-stack-4d3d5.firebaseapp.com",
  projectId: "credit-stack-4d3d5",
  storageBucket: "credit-stack-4d3d5.firebasestorage.app",
  messagingSenderId: "974632221642",
  appId: "1:974632221642:web:425a9c107a503862e22d33",
  measurementId: "G-LBWQ1TG0HK"
})

export const auth = getAuth(app);
export const firestore = getFirestore(app);