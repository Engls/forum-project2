import React from 'react';
import { BrowserRouter as Router, Route, Routes, Navigate } from 'react-router-dom';
import Login from './components/Auth/Login';
import Register from './components/Auth/Register';
import PostList from './components/Posts/PostList';
import CreatePost from './components/Posts/CreatePost';
import MainLayout from './components/Layout/MainLayout';
import styled from 'styled-components';

const AppContainer = styled.div`
  text-align: center;
`;

const PrivateRoute = ({ children }) => {
    const token = localStorage.getItem('token');
    return token ? children : <Navigate to="/login" />;
};

const App = () => {
    const onPostCreated = () => {
        console.log('Post was created');
    };

    return (
        <AppContainer>
            <Router>
                <MainLayout>
                    <Routes>
                        <Route path="/login" element={<Login />} />
                        <Route path="/register" element={<Register />} />
                        <Route path="/posts" element={
                            <PrivateRoute>
                                <CreatePost onPostCreated={onPostCreated} />
                                <PostList />
                            </PrivateRoute>
                        } />
                        <Route path="/" element={<Navigate to="/login" />} />
                    </Routes>
                </MainLayout>
            </Router>
        </AppContainer>
    );
};

export default App;