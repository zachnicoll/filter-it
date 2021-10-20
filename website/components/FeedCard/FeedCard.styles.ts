import styled from "@emotion/styled";
import theme, { colors } from "styles/colors";

export const Container = styled.div`
  display: flex;
  flex-direction: column;
  width: 100%;
  max-width: 600px;
  max-height: 800px;
  aspect-ratio: 0.9;
  background-color: ${colors.white};
  margin-bottom: 2em;
  border: 1px solid ${theme.border};
`;

export const ImageContainer = styled.div`
  width: 100%;
  aspect-ratio: 1;
  overflow: hidden;
`;

export const Info = styled.div`
  padding: 1em;
`;

export const Title = styled.p`
  font-weight: bold;
`;

export const Author = styled.p`
  font-size: 0.8em;
  font-style: italic;
  margin-top: 0.25em;
`;

export const Tag = styled.p`
  color: ${colors.hotBlue};
  font-size: 12px;
  margin-top: 0.5em;
`;
