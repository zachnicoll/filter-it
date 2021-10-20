import styled from "@emotion/styled";
import { Search } from "@material-ui/icons";
import theme, { colors } from "styles/colors";

const HEADER_HEIGHT = "3em";

export const Container = styled.main`
  display: flex;
  flex: 1;
  flex-direction: column;
  padding-top: calc(${HEADER_HEIGHT} + 2em);
  background-color: ${theme.background};
  overflow-y: auto;
  align-items: center;
  min-height: 100vh;
`;

export const HeaderContainer = styled.div`
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: center;

  position: fixed;
  top: 0;
  left: 0;

  background-color: ${colors.white};
  border-bottom: ${theme.border} 1px solid;

  height: ${HEADER_HEIGHT};
  width: 100%;
  padding-left: 1em;
  z-index: 1000;
`;

export const LinkContainer = styled.div`
  display: flex;
  flex-direction: row;
  align-items: center;
  flex: 1;
  height: 100%;
`;

export const LinkButton = styled.div<{ selected: boolean }>`
  padding: 0.5em;
  flex: 1;
  justify-content: center;
  align-items: center;
  display: flex;

  height: 100%;

  background-color: ${({ selected }) =>
    selected ? theme.primary : colors.white};
  color: ${({ selected }) => (selected ? colors.white : theme.text)};
`;

export const Title = styled.h1`
  flex: 2;
  font-style: bold;
  font-family: "Rouge Script", cursive;
  font-size: 36px;
`;

export const SearchInput = styled.input`
  flex: 2;
  margin: 0 0.5em 0 0.5em;
  border: 1px solid ${theme.border};
  border-radius: 2em;
  padding: 0.5em 1em;
`;

export const Mobile = {
  Header: styled(HeaderContainer)`
    padding: 0em 0.25em;
  `,
  LinkButton: styled(LinkButton)`
    height: fit-content;
  `,
};
