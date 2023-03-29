import { Form, Select } from 'antd';

type IProps = {};

const Index = (props: IProps) => {
  return (
    <>
      <Form.Item label={'ebpf lang'} name="ebpf_lang">
        <Select
          style={{
            width: 166,
          }}
        >
          <Select.Option value={0}>0</Select.Option>
          <Select.Option value={1}>1</Select.Option>
        </Select>
      </Form.Item>
    </>
  );
};
export default Index;
