// Mocks only back the official example pages; real APIs (/api/user/*, /api/system/*)
// are never registered here, so mockjs passes them through to the Go server.
export default ({ mock, setup }: { mock?: boolean; setup: () => void }) => {
  if (mock !== false) setup();
};

export const successResponseWrap = (data: unknown) => {
  return {
    data,
    status: 'ok',
    msg: '请求成功',
    code: 20000,
  };
};

export const failResponseWrap = (data: unknown, msg: string, code = 50000) => {
  return {
    data,
    status: 'fail',
    msg,
    code,
  };
};
