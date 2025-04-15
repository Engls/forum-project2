import React from 'react';
import styled from 'styled-components';
import Navbar from './Navbar';

const Content = styled.div`
  padding: 20px;
`;

const MainLayout = ({ children }) => {
    return (
        <>
            <Navbar />
            <Content>{children}</Content>
        </>
    );
};

export default MainLayout;