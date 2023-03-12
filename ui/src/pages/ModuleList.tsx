import { codeDelete, codeDeploy, codeList, codeUndeploy } from '@/services/ant-design-pro/api';
import { PageContainer } from '@ant-design/pro-components';
import { useIntl } from '@umijs/max';
import { Card, Form, message, Popconfirm, Table } from 'antd';
import type { ColumnsType } from 'antd/lib/table';
import React, { useEffect, useState } from 'react';

export interface CodeListItemType {
  Name: string;
  Status: string;
  CreateTime: number;
}

// maybe use
// const JsonRender = (json: any) => {
//   if (!json) {
//     return json;
//   }
//   if (typeof json === 'object') {
//     return <pre>{JSON.stringify(json, undefined, 4)}</pre>;
//   }
//   return json;
// };

const ArrayRender = (
  data: {
    label: string;
    value: any;
  }[],
  layout = 'horizontal',
) => {
  if (!Array.isArray(data)) {
    return data;
  }
  return (
    <Form layout={layout as any}>
      {data.map((item) => {
        if (Array.isArray(item)) {
          return ArrayRender(item, 'inline');
        }
        return (
          <Form.Item label={item.label} key={item.value}>
            {item.value}
          </Form.Item>
        );
      })}
    </Form>
  );
};

const CodeList: React.FC = () => {
  const [data, setData] = useState<CodeListItemType[]>();
  const intl = useIntl();
  const getData = async () => {
    try {
      const msg: any = await codeList({
        fields: `""`,
      });
      console.log('codeList', msg);
      if (msg.code != 200) {
        message.error(msg.message);
        return;
      }
      setData(msg.data);
      return msg.data;
    } catch (error) {
      console.log(error);
      message.error('请求失败');
      return error;
    }
  };
  const columns: ColumnsType<CodeListItemType> = [
    {
      title: intl.formatMessage({
        id: 'codeList.name',
      }),
      dataIndex: 'name',
      render: (val) => {
        return val;
      },
    },
    {
      title: intl.formatMessage({
        id: 'codeList.status',
      }),
      // 0: 代表未部署，1代表成功 2 部署中 3 部署失败
      dataIndex: 'status',
      render: (d: number) => {
        if (d === 0) {
          return intl.formatMessage({
            id: 'codeList.status.undeploy',
          });
        }
        if (d === 1) {
          return intl.formatMessage({
            id: 'codeList.status.success',
          });
        }
        if (d === 2) {
          return intl.formatMessage({
            id: 'codeList.status.pending',
          });
        }
        if (d === 3) {
          return intl.formatMessage({
            id: 'codeList.status.fail',
          });
        }
        return '--';
      },
    },
    {
      title: intl.formatMessage({
        id: 'codeList.create_time',
      }),
      dataIndex: 'create_time',
    },
    {
      title: intl.formatMessage({
        id: 'codeList.action',
      }),
      dataIndex: '',
      key: 'x',
      render: (_: any, columnData: any) => (
        <>
          <a
            onClick={async () => {
              const res = await codeDeploy({
                Id: columnData.id || columnData.ID,
              });
              if (res.code === 200) {
                message.success('deploy success');
                getData();
              } else {
                message.error('deploy fail' + res.message);
              }
            }}
          >
            {intl.formatMessage({
              id: 'codeList.action.deploy',
            })}
          </a>
          <a
            onClick={async () => {
              const res = await codeUndeploy({
                Id: columnData.id || columnData.ID,
              });
              if (res.code === 200) {
                message.success('undeploy success');
                getData();
              } else {
                message.error('undeploy fail,' + res.message);
              }
            }}
            style={{
              marginLeft: 10,
            }}
          >
            {intl.formatMessage({
              id: 'codeList.action.undeploy',
            })}
          </a>

          <Popconfirm
            placement="bottom"
            title={'Are you sure to delete ?'}
            onConfirm={async () => {
              const res = await codeDelete({
                Id: columnData.id || columnData.ID,
              });
              debugger;
              if (res.code === 200) {
                message.success('delete success');
                getData();
              } else {
                message.error('delete fail, ' + res.message);
              }
            }}
            okText="Yes"
            cancelText="No"
          >
            <a
              style={{
                marginLeft: 10,
              }}
            >
              {intl.formatMessage({
                id: 'codeList.action.delete',
              })}
            </a>
          </Popconfirm>
        </>
      ),
    },
  ];
  useEffect(() => {
    getData();
  }, []);
  return (
    <PageContainer>
      <Card
        style={{
          borderRadius: 8,
        }}
        bodyStyle={{
          backgroundImage:
            'radial-gradient(circle at 97% 10%, #EBF2FF 0%, #F5F8FF 28%, #EBF1FF 124%)',
        }}
      >
        <Table
          columns={columns}
          expandable={{
            expandedRowRender: (columnData) => {
              return (
                <Table
                  style={{
                    padding: '30px 20px 30px 0',
                  }}
                  showHeader={false}
                  columns={[
                    {
                      title: 'id',
                      dataIndex: 'ID',
                    },
                    {
                      title: 'name',
                      dataIndex: 'Name',
                    },
                    {
                      title: 'Ebpf',
                      dataIndex: 'Ebpf',
                      render: (da: string) => {
                        if (!da || da === 'Ebpf') {
                          return da;
                        }
                        const res: any = JSON.parse(da || '{}');
                        return (
                          <Form>
                            <Form.Item label="code">{res.code}</Form.Item>
                            <Form.Item label="eventSize">{res.eventSize}</Form.Item>
                            <Form.Item label="perfBuffers">
                              {ArrayRender(
                                res.kprobes?.map((i: any) => [
                                  {
                                    label: 'target',
                                    value: i.target,
                                  },
                                  {
                                    label: 'entry',
                                    value: i.entry,
                                  },
                                  {
                                    label: 'return',
                                    value: i.return,
                                  },
                                ]),
                              )}
                            </Form.Item>
                          </Form>
                        );
                      },
                    },
                    {
                      title: 'SchemaAttr',
                      dataIndex: 'SchemaAttr',
                      render: (da: string) => {
                        if (!da || da === 'SchemaAttr') return da;
                        const res: any = JSON.parse(da || '[]');
                        return ArrayRender(
                          res?.map((i: any) => ({
                            label: i.name,
                            value: i.type,
                          })),
                        );
                      },
                    },
                    {
                      title: 'WasmFileName',
                      dataIndex: 'WasmFileName',
                    },
                  ]}
                  dataSource={[
                    {
                      ID: 'ID',
                      Name: 'name',
                      Ebpf: 'Ebpf',
                      SchemaAttr: 'SchemaAttr',
                      WasmFileName: 'WasmFileName',
                    },
                    columnData,
                  ]}
                  bordered
                />
              );
            },
          }}
          rowKey={(item: any) => item.ID}
          dataSource={data}
        />
      </Card>
    </PageContainer>
  );
};

export default CodeList;
