import { FeedCard, PageWrapper, Spinner } from "components";
import { useSearch } from "context";
import { useImages } from "hooks/useImages";
import type { NextPage } from "next";
import React from "react";

const Feed: NextPage = () => {
  const { searchState } = useSearch();
  const { images, loading } = useImages(searchState.search);

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
