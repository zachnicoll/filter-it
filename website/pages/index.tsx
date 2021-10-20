import { FeedCard, PageWrapper, Spinner } from "components";
import { useImages } from "hooks/useImages";
import type { NextPage } from "next";
import React from "react";

const Feed: NextPage = () => {
  const { images, loading } = useImages();

  return (
    <PageWrapper>
      {!loading ? (
        images.map((image) => <FeedCard key={image.id} image={image} />)
      ) : (
        <Spinner type="Oval" />
      )}
    </PageWrapper>
  );
};

export default Feed;
