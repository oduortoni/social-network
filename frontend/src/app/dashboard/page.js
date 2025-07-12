"use client";

// import { useEffect } from "react";
// import { useRouter } from "next/navigation";

export default function Dashboard() {
  // const router = useRouter();

  // useEffect(() => {
  //   fetch('/api/me', { credentials: 'include' })
  //     .then(res => res.json())
  //     .then(data => {
  //       console.log(data);
  //     })
  //     .catch(err => {
  //       console.log(err);
  //       router.push('/login');
  //     });
  // }, [router]);

  return (
    <div style={{ textAlign: 'center', marginTop: '50px' }}>
      <h1 style={{ fontSize: '30px', fontWeight: 'bold', color: '#fff' }}>
        Welcome to your Dashboard!
      </h1>
    </div>
  );  
}
