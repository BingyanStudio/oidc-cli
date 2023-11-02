import React, { useEffect, useState } from 'react';
import { UserManager, UserManagerSettings } from 'oidc-client-ts';

// OIDC client configuration
const oidcConfig: UserManagerSettings = {
    authority: 'https://api.bingyan.net/sso/oidc',
    client_id: 'CLIENT_ID',
    redirect_uri: 'http://localhost:3000/callback', // Frontend Callback URI
    response_type: 'code',
    scope: 'openid profile phone email',
};
const retrieve_token_uri = 'http://localhost:8000/user/token'; // Backend URI to retrieve the tokens

const userManager = new UserManager(oidcConfig);

const OidcClient: React.FC = () => {
    const [status, setStatus] = useState<boolean>(false);
    useEffect(() => {
        const exchangeCode = async () => {
            const params = new URLSearchParams(window.location.search);
            const code = params.get('code');
            const state = params.get('state');

            if (code && state) {
                const response = await fetch(`${retrieve_token_uri}?code=${code}&state=${state}`, {
                    method: 'GET',
                });
                const data = await response.json();
                console.log(data)
                setStatus(true)
            }
        };
        exchangeCode();
    }, []);

    const handleLogin = () => {
        userManager.signinRedirect();
    };

    const handleLogout = () => {
        userManager.signoutRedirect();
    };

    return (
        <div>
            {status ? (
                <div>
                    <div className="user-data">
                        <h2>Log in successful. View console for more info</h2>
                    </div>
                    <button onClick={handleLogout}>Logout</button>
                </div>
            ) : (
                <button onClick={handleLogin}>Login</button>
            )}
        </div>
    );
};

export default OidcClient;