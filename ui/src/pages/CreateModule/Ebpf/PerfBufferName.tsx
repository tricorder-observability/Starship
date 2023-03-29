import { Form, Input } from 'antd';

type IProps = {};

const Index = (props: IProps) => {
  return (
    <Form.Item
      label={'perf buffer name'}
      name="perf_buffer_name"
      rules={[{ required: true, message: 'Please input eBPF perf buffers!' }]}
    >
      <Input style={{ width: '100%' }} />
    </Form.Item>
  );
};
export default Index;
