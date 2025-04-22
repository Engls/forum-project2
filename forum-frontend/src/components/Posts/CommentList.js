import React, { useState, useEffect } from 'react';
import axios from 'axios';
import styled from 'styled-components';

const CommentListContainer = styled.div`
    margin-top: 15px;
    padding: 10px;
    border-radius: 8px;
    background-color: #f8f8f8;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
`;

const CommentItem = styled.div`
    padding: 8px;
    margin-bottom: 8px;
    border-bottom: 1px solid #eee;
    &:last-child {
        border-bottom: none;
    }
`;

const CommentContent = styled.p`
    font-size: 14px;
    color: #333;
    margin-bottom: 5px;
`;

const CommentAuthor = styled.small`
    color: #777;
    font-style: italic;
`;

const LoadingMessage = styled.p`
    color: #555;
    font-style: italic;
`;

const ErrorMessage = styled.p`
    color: #d32f2f;
`;

const NoCommentsMessage = styled.p`
    color: #999;
    font-style: italic;
`;

const CommentList = ({ postId }) => {
    const [comments, setComments] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        const fetchComments = async () => {
            setLoading(true);
            setError(null);
            try {
                const response = await axios.get(`http://localhost:8081/posts/${Number(postId)}/comments`);
                setComments(response.data ?? []);
            } catch (error) {
                console.error('Error fetching comments:', error);
                setError(error);
            } finally {
                setLoading(false);
            }
        };

        fetchComments();
    }, [postId]);

    if (loading) return <LoadingMessage>Loading comments...</LoadingMessage>;
    if (error) return <ErrorMessage>Error: {error.message} {error.response?.status && `(Status: ${error.response.status})`}</ErrorMessage>;

    return (
        <CommentListContainer>
            <h4>Comments:</h4>
            {comments.length ? (
                comments.map(comment => (
                    <CommentItem key={comment.id}>
                        <CommentContent>{comment.content}</CommentContent>
                        <CommentAuthor>By User ID: {comment.author_id}</CommentAuthor>
                    </CommentItem>
                ))
            ) : (
                <NoCommentsMessage>No comments yet.</NoCommentsMessage>
            )}
        </CommentListContainer>
    );
};

export default CommentList;