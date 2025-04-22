import React, { useEffect, useState, useCallback } from 'react';
import axios from 'axios';
import styled from 'styled-components';
import CommentList from './CommentList';
import AddComment from './AddComment';

// === Styled Components ===

const Container = styled.div`
    max-width: 800px;
    margin: 20px auto;
    padding: 20px;
    background-color: #fff;
    border-radius: 8px;
    box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
`;

const PostItemContainer = styled.div`
    margin-bottom: 20px;
    padding: 15px;
    border: 1px solid #e0e0e0;
    border-radius: 6px;
    background-color: #f9f9f9;
`;

const PostTitle = styled.h3`
    color: #333;
    text-align: center; /* Заголовок по центру */
    margin-bottom: 5px;
    font-size:32px;
`;

const PostContent = styled.p`
    color: #555;
    font-size: 16px;
    margin-bottom: 10px;
    text-align: left; /* Текст по левому краю */
`;

const PostAuthor = styled.small`
    color: #777;
    font-style: italic;
`;

const DeleteButton = styled.button`
    background-color: #f44336;
    color: white;
    padding: 8px 12px;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    transition: background-color 0.2s ease;

    &:hover {
        background-color: #d32f2f;
    }
    margin-top: 10px;
`;

const LoadingMessage = styled.p`
    color: #555;
    font-style: italic;
    text-align: center;
`;

const ErrorMessage = styled.p`
    color: #d32f2f;
    text-align: center;
`;

// === PostList Component ===

const PostList = () => {
    const [posts, setPosts] = useState([]);
    const [isAdmin, setIsAdmin] = useState(false);
    const [forceUpdate, setForceUpdate] = useState(false);

    const fetchPosts = useCallback(async () => {
        try {
            const response = await axios.get('http://localhost:8081/posts');
            setPosts(response.data);
        } catch (error) {
            console.error('Error fetching posts:', error);
        }
    }, []);

    useEffect(() => {
        const storedRole = localStorage.getItem('role');
        setIsAdmin(storedRole === 'admin');
        fetchPosts();
        const intervalId = setInterval(fetchPosts, 5000);
        return () => clearInterval(intervalId);
    }, [fetchPosts]);

    const handleDeletePost = async (postId) => {
        const token = localStorage.getItem('token');

        if (!token) {
            alert('You are not authenticated.');
            return;
        }

        try {
            await axios.delete(`http://localhost:8081/posts/${postId}`, {
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            });
            setPosts(posts.filter(post => post.id !== postId));
        } catch (error) {
            console.error('Error deleting post:', error);
            if (error.response && error.response.status === 403) {
                alert('You are not authorized to delete this post.');
            } else {
                alert('Failed to delete post.');
            }
        }
    };

    const handleCommentCreated = () => {
        setForceUpdate(prev => !prev);
    };
    if(posts===null){
        return <LoadingMessage>No posts available.</LoadingMessage>;
    }
    if (posts.length === 0 && !isAdmin) {
        return <LoadingMessage>No posts available.</LoadingMessage>;
    }

    return (
        <Container>
            {posts.length === 0 && isAdmin ? (
                <LoadingMessage>No posts yet.</LoadingMessage>
            ) : (
            [...posts].reverse().map(post => (  // Разворачиваем массив здесь
                <PostItemContainer key={post.id}>
                    <PostTitle>{post.title}</PostTitle>
                    <PostContent>{post.content}</PostContent>
                    <PostAuthor>Created by User ID: {post.author_id}</PostAuthor>
                    <CommentList postId={post.id} />
                    <AddComment postId={post.id} onCommentCreated={handleCommentCreated} />
                    {isAdmin && (
                        <DeleteButton onClick={() => handleDeletePost(post.id)}>
                            Delete
                        </DeleteButton>
                    )}
                </PostItemContainer>
            ))
            )}
        </Container>
    );
};

export default PostList;