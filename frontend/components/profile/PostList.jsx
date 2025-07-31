
const PostList = () => {
  const posts = [1, 2, 3, 4, 5, 6];

  return (
    <div className="mt-8">
      <h2 className="text-xl font-bold mb-4">Posts</h2>
      <div className="grid grid-cols-3 gap-4">
        {posts.map((post) => (
          <div key={post} className="bg-gray-200 h-40 rounded-lg"></div>
        ))}
      </div>
    </div>
  );
};

export default PostList;
