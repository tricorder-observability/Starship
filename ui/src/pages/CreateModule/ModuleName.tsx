import { useIntl } from '@umijs/max';
import { Form, Input } from 'antd';

type IProps = {};

const Index = (props: IProps) => {
  const intl = useIntl();
  return (
    <Form.Item
      label={intl.formatMessage({
        id: 'code.name',
      })}
      name="name"
      rules={[
        { required: true, message: 'Please input name!' },
        { type: 'string', max: 50, message: 'Up to 50 characters!' },
      ]}
    >
      <Input style={{ width: '100%' }} />
    </Form.Item>
  );
};
export default Index;
