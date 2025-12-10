const sleep = (ms) => new Promise((resolve) => setTimeout(resolve, ms));

class App {
    // 请实现，不限制一定要使用 class 实现，也可以用原型链或其它方法
    fnList = [];
    ctx = {};
    push(v) {
        this.ctx.push.push(v);
    }
    use(fn) {
        this.fnList.push(fn);
    }
    async run(ctx) {
        this.ctx = ctx;
        const fn = this.fnList.reduceRight(((curr, prev) => {
            return prev.bind(this, this, curr);
        }), sleep);
        await fn();
    }
}

// --------------------以下为执行示例-----------------------

(async () => {
    const app = new App();

    app.use(async (ctx, next) => {
        ctx.push(1);
        await next();
        ctx.push(4);
    });

    app.use(async (ctx, next) => {
        ctx.push(2);
        await next();
        ctx.push(3);
    });

    const ctx = { push: [] };
    await app.run(ctx);
    console.log(ctx.push); // 期望输出：[1, 2, 3, 4]
})();

