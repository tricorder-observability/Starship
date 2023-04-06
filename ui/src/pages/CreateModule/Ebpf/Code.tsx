import { Form, Input } from 'antd';

type IProps = {};

const Index = (props: IProps) => {
  return (
    <>
      <Form.Item
        label={'ebpf code'}
        name="ebpf_code"
        rules={[{ required: true, message: 'Please input ebpf!' }]}
      >
        <Input style={{ width: '100%' }} />
      </Form.Item>
    </>
  );
};
export default Index;
