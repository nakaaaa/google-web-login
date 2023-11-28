"use client"
import { useEffect } from "react";
import { useSearchParams } from "next/navigation";

type TokenInfo = {
  AccessToken: string;
  ExpiresIn: number;
  IDToken: string;
  Scope: string;
  TokenType: string;
  RefreshToken: string;
}

type IDToken = {
  Iss: string;
  Azp: string;
  Aud: string;
  Sub: string;
  AtHash: string;
  HD: string;
  Email: string;
  EmailVerified: string;
  Iat: string;
  Exp: string;
  Nonce: string;
}

const CallbackPage = () => {
  const params = useSearchParams();

  useEffect(() => {
    const code = params.get('code');

    if (code) {
      fetch(`http://localhost:8080/verify/id_token?code=${code}`)
        .then((response) => response.json())
        .then((data: IDToken) => {
          console.log(data);
        })
        .catch((error) => {
          // エラーハンドリングを行う
        });
    } else {
      // codeパラメータが存在しない場合の処理を行う
    }
  }, [params]);

  return (
    <div>
      CallbackPage
    </div>
  );
};

export default CallbackPage;