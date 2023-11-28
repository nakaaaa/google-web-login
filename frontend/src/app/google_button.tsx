import React from 'react';

type Auth = {
  url: string;
}

const GoogleButton: React.FC = () => {
  const handleClick = () => {
    fetch('http://localhost:8080/auth', {
      method: 'GET',
    })
      .then(response =>response.json())
      .then((data: Auth) => {
        window.location.href = data.url;
      })
      .catch(error => {
        console.error(error);
      });
  };

  return (
    <button onClick={handleClick}>Google ログイン</button>
  );
};

export default GoogleButton;
