
'use client';
import { useState } from 'react';
import ProfileHeader from "../../../components/profile/ProfileHeader";
import PostList from "../../../components/profile/PostList";
import FollowerList from "../../../components/profile/FollowerList";
import FollowingList from "../../../components/profile/FollowingList";

const ProfilePage = () => {
  const [isPublic, setIsPublic] = useState(true);

  return (
    <div className="container mx-auto p-4">
      <div className="flex justify-end mb-4">
        <button onClick={() => setIsPublic(!isPublic)} className="bg-gray-200 px-4 py-2 rounded-lg">
          Toggle Private/Public
        </button>
      </div>
      <ProfileHeader />
      {isPublic ? (
        <>
          <PostList />
          <FollowerList />
          <FollowingList />
        </>
      ) : (
        <p className="text-center mt-8">This account is private.</p>
      )}
    </div>
  );
};

export default ProfilePage;
