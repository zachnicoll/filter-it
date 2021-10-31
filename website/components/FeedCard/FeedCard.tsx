import { ImageDocument } from "api/types";
import React from "react";
import * as Styles from "./FeedCard.styles";
import { filterStringMap } from "common/enum.maps";

interface FeedCardProps {
  image: ImageDocument;
}

export const FeedCard: React.FC<FeedCardProps> = ({ image }) => {
  return (
    <Styles.Container>
      <Styles.ImageContainer>
        <img
          src={decodeURIComponent(image.image)}
          width="100%"
          height="100%"
          alt="image"
        />
      </Styles.ImageContainer>

      <Styles.Info>
        <Styles.Title>{image.title}</Styles.Title>
        <Styles.Author>{image.author}</Styles.Author>

        <Styles.Tag>{filterStringMap[image.tag]}</Styles.Tag>
      </Styles.Info>
    </Styles.Container>
  );
};
