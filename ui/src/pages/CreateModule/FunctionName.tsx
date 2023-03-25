import { useIntl } from '@umijs/max';
import { Form, Input } from 'antd';

type IProps = {};

const Index = (props: IProps) => {
  const intl = useIntl();
  return (
    <Form.Item
      label={intl.formatMessage({
        id: 'code.fn',
      })}
      name="fn"
      rules={[{ required: true, message: 'Please input function name!' }]}
    >
      <Input style={{ width: '100%' }} />
    </Form.Item>
  );
};
export default Index;
