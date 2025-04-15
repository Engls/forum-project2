import React, { useEffect, useState } from 'react';
import axios from 'axios';
import styled from 'styled-components';

const PostListContainer = styled.div`
  margin-top: 20px;
`;

const PostItem = styled.div`
  padding: 15px;
  border: 1px solid #ddd;
  border-radius: 5px;
  margin-bottom: 10px;
  background-color: #f9f9f9;
`;

const PostList = () => {
    const [posts, setPosts] = useState([]);

    useEffect(() => {
        const fetchPosts = async () => {
            try {
                const response = await axios.get('http://localhost:8081/posts');
                setPosts(response.data);
            } catch (error) {
                console.error('Error fetching posts:', error);
            }
        };

        fetchPosts();
    }, []);

    return (
        <PostListContainer>
            {posts && posts.length > 0 ? (
                [...posts].reverse().map(post => (
                    <PostItem key={post.id}>
                        <h3>{post.title}</h3>
                        <p>{post.content}</p>
                        <small>Created by User ID: {post.author_id}</small>
                    </PostItem>
                ))
            ) : (
                <p>No posts available.</p>
            )}
        </PostListContainer>
    );
};

export default PostList;