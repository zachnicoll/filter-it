import styled from "@emotion/styled";
import theme from "../../styles/colors";

export const UploadWrapper = styled.div`
  width: 100%;
  max-width: 600px;
  padding: 10px;
  background-color: ${theme.white};
  border-width: 2px;
  border-radius: 1px;
  border-color: ${theme.border};
  border-style: solid;
  align-self: center;
`;

export const Img = styled.img`
  max-width: 100%;
  height: auto;
  display: block;
  margin: auto;
`;

export const Container = styled.div<{ active?: boolean }>`
  align-items: center;
  border-width: 2px;
  border-radius: 2px;
  background-color: ${({ active }) => (active ? theme.accent : theme.white)};
  border-style: dashed;
  padding: 10px;
  color: ${theme.secondary};
  outline: none;
  text-align: center;
`;

export const FormColumn = styled.div`
  display: flex;
  flex-direction: column;
  margin-right: 10%;
`;

export const FormRow = styled.div`
  display: flex;
  flex-direction: row;
  margin-bottom: 5px;
  flex: 1;
  align-items: center;
  flex-wrap: wrap;
`;

export const FormArea = styled.div`
  margin: 10px 0px;
  display: flex;
  flex-direction: column;
  flex: 1;
`;

export const FormBox = styled.div`
  display: flex;
  flex-direction: column;
  margin: 5px;
  flex: 1;
`;

export const FormLabel = styled.label<{ emphasized?: boolean }>`
  font-weight: ${({ emphasized }) => (emphasized ? `bold` : `normal`)};
  color: ${theme.secondary};
  text-transform: capitalize;
  margin-bottom: 5px;
`;

export const FormCheckBoxInput = styled.input`
  margin-right: 5px;
`;

export const FormButton = styled.button`
  border: 0;
  border-radius: 0.25rem;
  font-size: 1rem;
  line-height: 1.2;
  white-space: nowrap;
  text-decoration: none;
  padding: 0.5rem;
  cursor: pointer;
  color: ${theme.primary};

  &:hover {
    background: ${theme.primary};
    color: ${theme.white};
  }
`;

export const Input = styled.input`
  border: 1px solid ${theme.border};
  border-radius: 0.5em;
  padding: 0.25em 0.5em;
`;
