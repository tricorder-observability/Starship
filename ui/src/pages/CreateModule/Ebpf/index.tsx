import Code from './Code';
import Fmt from './Fmt';
import Lang from './Lang';
import PerfBufferName from './PerfBufferName';
import Probes from './Probes';

type IProps = {};

const Index = (props: IProps) => {
  return (
    <>
      <Code />
      <Fmt />
      <Lang />
      <PerfBufferName />
      <Probes />
    </>
  );
};
export default Index;
