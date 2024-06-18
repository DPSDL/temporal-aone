import React, { useState } from "react";

const App = () => {
    const [repoURL, setRepoURL] = useState("");
    const [token, setToken] = useState("");
    const [binaryPath, setBinaryPath] = useState("");
    const [configFilePath, setConfigFilePath] = useState("");
    const [version, setVersion] = useState("");
    const [ecsUploadPath, setECSUploadPath] = useState("");
    const [ecsServer, setECSServer] = useState("");
    const [ecsUser, setECSUser] = useState("");
    const [ecsPassword, setECSPassword] = useState("");
    const [healthCheckURL, setHealthCheckURL] = useState("");

    const handleStartConfig = async () => {
        const response = await fetch("/api/start-config", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                repo_url: repoURL,
                token: token,
                binary_path: binaryPath,
                config_file_path: configFilePath,
                version: version,
                ecs_upload_path: ecsUploadPath,
                ecs_server: ecsServer,
                ecs_user: ecsUser,
                ecs_password: ecsPassword,
                health_check_url: healthCheckURL,
            }),
        });

        const data = await response.json();
        console.log(data);
    };

    const handleStartBuildUpload = async () => {
        const response = await fetch("/api/start-build-upload", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                repo_url: repoURL,
                token: token,
                binary_path: binaryPath,
                config_file_path: configFilePath,
                version: version,
                ecs_upload_path: ecsUploadPath,
                ecs_server: ecsServer,
                ecs_user: ecsUser,
                ecs_password: ecsPassword,
                health_check_url: healthCheckURL,
            }),
        });

        const data = await response.json();
        console.log(data);
    };

    const handleStartRelease = async () => {
        const response = await fetch("/api/start-release", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                repo_url: repoURL,
                token: token,
                binary_path: binaryPath,
                config_file_path: configFilePath,
                version: version,
                ecs_upload_path: ecsUploadPath,
                ecs_server: ecsServer,
                ecs_user: ecsUser,
                ecs_password: ecsPassword,
                health_check_url: healthCheckURL,
            }),
        });

        const data = await response.json();
        console.log(data);
    };

    return (
        <div>
            <h1>Workflow Config</h1>
            <div>
                <label>Repo URL:</label>
                <input
                    type="text"
                    value={repoURL}
                    onChange={(e) => setRepoURL(e.target.value)}
                />
            </div>
            <div>
                <label>Token:</label>
                <input
                    type="text"
                    value={token}
                    onChange={(e) => setToken(e.target.value)}
                />
            </div>
            <div>
                <label>Binary Path:</label>
                <input
                    type="text"
                    value={binaryPath}
                    onChange={(e) => setBinaryPath(e.target.value)}
                />
            </div>
            <div>
                <label>Config File Path:</label>
                <input
                    type="text"
                    value={configFilePath}
                    onChange={(e) => setConfigFilePath(e.target.value)}
                />
            </div>
            <div>
                <label>Version:</label>
                <input
                    type="text"
                    value={version}
                    onChange={(e) => setVersion(e.target.value)}
                />
            </div>
            <div>
                <label>ECS Upload Path:</label>
                <input
                    type="text"
                    value={ecsUploadPath}
                    onChange={(e) => setECSUploadPath(e.target.value)}
                />
            </div>
            <div>
                <label>ECS Server:</label>
                <input
                    type="text"
                    value={ecsServer}
                    onChange={(e) => setECSServer(e.target.value)}
                />
            </div>
            <div>
                <label>ECS User:</label>
                <input
                    type="text"
                    value={ecsUser}
                    onChange={(e) => setECSUser(e.target.value)}
                />
            </div>
            <div>
                <label>ECS Password:</label>
                <input
                    type="text"
                    value={ecsPassword}
                    onChange={(e) => setECSPassword(e.target.value)}
                />
            </div>
            <div>
                <label>Health Check URL:</label>
                <input
                    type="text"
                    value={healthCheckURL}
                    onChange={(e) => setHealthCheckURL(e.target.value)}
                />
            </div>
            <div>
                <button onClick={handleStartConfig}>Start Config Workflow</button>
                <button onClick={handleStartBuildUpload}>Start Build & Upload Workflow</button>
                <button onClick={handleStartRelease}>Start Release Workflow</button>
            </div>
        </div>
    );
};

export default App;
