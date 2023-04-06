import Code from './Code';
import Fmt from './Fmt';
import FunctionName from './FunctionName';
import Lang from './Lang';
import OutputSchema from './OutputSchema';

type IProps = {
  readFileContent: (info: any) => void;
};

const Index = (props: IProps) => {
  const { readFileContent } = props;
  return (
    <>
      <Code readFileContent={readFileContent} />
      <Fmt />
      <FunctionName />
      <Lang />
      <OutputSchema />
    </>
  );
};
export default Index;
